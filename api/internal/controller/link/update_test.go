package link

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

func TestUpdateURL(t *testing.T) {
	userID := uuid.New()
	const (
		code       = "abc1234"
		rawURL     = "https://Example.com/Path"
		normalized = "https://example.com/Path"
	)

	tcs := map[string]struct {
		setup         func(*repolink.MockRepository)
		wantURL       string
		wantErrStatus int
	}{
		"updates the destination": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().UpdateURL(mock.Anything, code, userID, rawURL, normalized).
					Return(model.Link{Code: code, OriginalURL: rawURL}, nil).Once()
			},
			wantURL: rawURL,
		},
		"conflict maps to 409": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().UpdateURL(mock.Anything, code, userID, rawURL, normalized).
					Return(model.Link{}, repolink.ErrConflict).Once()
			},
			wantErrStatus: http.StatusConflict,
		},
		"not found maps to 404": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().UpdateURL(mock.Anything, code, userID, rawURL, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
			},
			wantErrStatus: http.StatusNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			repo := repolink.NewMockRepository(t)
			tc.setup(repo)

			link, err := New(repo).UpdateURL(context.Background(), code, userID, rawURL)
			if tc.wantErrStatus != 0 {
				var ae *apperror.Error
				require.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantErrStatus, ae.HTTPStatus())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantURL, link.OriginalURL)
		})
	}
}
