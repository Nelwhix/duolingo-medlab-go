package models

import (
	"context"
	"time"
)

func (m *Model) InsertIntoPasswordResets(ctx context.Context, email, tokenHash string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := m.DeletePasswordReset(ctx, email)
	if err != nil {
		return err
	}

	query := "insert into password_reset_tokens (email, token_hash, expires_at) values ($1, $2, $3)"
	_, err = m.Conn.Exec(ctx, query, email, tokenHash, expiresAt)

	return err
}

func (m *Model) DeletePasswordReset(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "delete from password_reset_tokens where email = $1"
	_, err := m.Conn.Exec(ctx, query, email)

	return err
}

func (m *Model) GetPasswordReset(ctx context.Context, email string) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var tokenHash string
	var expiresAt time.Time
	query := "select token_hash, expires_at from password_reset_tokens where email = $1"
	err := m.Conn.QueryRow(ctx, query, email).Scan(&tokenHash, &expiresAt)

	return tokenHash, expiresAt, err
}

func (m *Model) UpdateUserPassword(ctx context.Context, email, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "update users set password = $1 where email = $2"
	_, err := m.Conn.Exec(ctx, query, passwordHash, email)

	return err
}
