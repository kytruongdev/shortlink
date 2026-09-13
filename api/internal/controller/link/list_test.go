package link

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

func TestListByUser(t *testing.T) {
	userID := uuid.New()
	want := []model.Link{{Code: "b"}, {Code: "a"}}

	repo := repolink.NewMockRepository(t)
	repo.EXPECT().ListByUserID(mock.Anything, userID).Return(want, nil).Once()

	got, err := New(repo).ListByUser(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
