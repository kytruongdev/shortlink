package authctx

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserContext(t *testing.T) {
	id := uuid.New()

	got, ok := UserFromContext(WithUser(context.Background(), id))
	require.True(t, ok)
	assert.Equal(t, id, got)

	_, ok = UserFromContext(context.Background())
	assert.False(t, ok)
}
