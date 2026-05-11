package context_key

import (
	"context"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(models.User)
	return user, ok
}
