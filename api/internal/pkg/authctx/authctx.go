package authctx

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey struct{}

// WithUser returns a copy of ctx carrying the authenticated user id.
func WithUser(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, userID)
}

// UserFromContext returns the authenticated user id, or false for an anonymous request.
func UserFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKey{}).(uuid.UUID)
	return id, ok
}
