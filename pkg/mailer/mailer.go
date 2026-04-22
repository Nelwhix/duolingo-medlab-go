package mailer

import "context"

type Mailer interface {
	SendPasswordResetEmail(ctx context.Context, to string, token string) error
}
