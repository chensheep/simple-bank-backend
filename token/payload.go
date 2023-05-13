package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("token is invalid")
var ErrExpiredToken = errors.New("token has expired")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Payload{
		ID:        uuid,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}, err
}

func (p *Payload) Valid() error {
	if p.ExpiredAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
