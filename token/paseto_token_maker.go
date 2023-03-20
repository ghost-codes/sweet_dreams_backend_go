package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Error occured:invalid key size; Size must be %d", chacha20poly1305.KeySize)

	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}

	return maker, nil
}

func (paseto *PasetoMaker) CreateToken(userId int64, securityKey string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, securityKey, duration)

	if err != nil {
		return "", payload, nil
	}

	token, err := paseto.paseto.Encrypt(paseto.symmetricKey, payload, nil)

	return token, payload, err
}

func (paseto *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := paseto.paseto.Decrypt(token, paseto.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, err
}
