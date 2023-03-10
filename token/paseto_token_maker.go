package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func (paseto *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", payload, nil
	}

	token, err := paseto.paseto.Encrypt(paseto.symmetricKey, payload, nil)

	return token, payload, err
}

func (paseto *PasetoMaker) ValidateToken(token string) (*Payload, error) {
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
