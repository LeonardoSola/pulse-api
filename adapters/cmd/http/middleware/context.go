package middleware

import "context"

type contextKey string

const authUserIDKey contextKey = "userID"

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(authUserIDKey).(string)
	return id, ok
}
