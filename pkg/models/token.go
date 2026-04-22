package models

import (
	"context"
	"time"
)

type CreateTokenRequest struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func (m *Model) InsertIntoTokens(ctx context.Context, request CreateTokenRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "insert into personal_access_tokens (user_id, token, expires_at) values ($1, $2, $3)"
	_, err := m.Conn.Exec(ctx, query, request.UserID, request.Token, request.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}
