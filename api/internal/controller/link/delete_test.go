package link

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

func TestDelete(t *testing.T) {
	userID := uuid.New()
	const code = "abc1234"

	tcs := map[string]struct {
		setup         func(*repolink.MockRepository)
		wantErrStatus int
	}{
		"deletes": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().Delete(mock.Anything, code, userID).Return(nil).Once()
			},
		},
		"not found maps to 404": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().Delete(mock.Anything, code, userID).Return(repolink.ErrNotFound).Once()
			},
			wantErrStatus: http.StatusNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			repo := repolink.NewMockRepository(t)
			tc.setup(repo)

			err := New(repo).Delete(context.Background(), code, userID)
			if tc.wantErrStatus != 0 {
				var ae *apperror.Error
				require.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantErrStatus, ae.HTTPStatus())
				return
			}
			require.NoError(t, err)
		})
	}
}
