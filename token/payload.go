package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorInvalidToken = errors.New("token is invalid")
	ErrorExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	SecurityKey string    `json:"security_key"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func NewPayload(username string, securityKey string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenId,
		Username:    username,
		IssuedAt:    time.Now(),
		SecurityKey: securityKey,
		ExpiredAt:   time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}

	return nil
}
