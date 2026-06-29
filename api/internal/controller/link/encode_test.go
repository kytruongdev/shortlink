package link

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

func TestEncode(t *testing.T) {
	const rawURL = "https://Example.com/Path"
	const normalized = "https://example.com/Path"

	tcs := map[string]struct {
		url           string
		setup         func(*repolink.MockRepository)
		wantCode      string // "" = random code, only length is checked
		wantErrStatus int    // 0 = no error
	}{
		"dedup returns existing": {
			url: rawURL,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{Code: "OLDCODE"}, nil).Once()
			},
			wantCode: "OLDCODE",
		},
		"new url creates code": {
			url: rawURL,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, l model.Link) (model.Link, error) { return l, nil }).Once()
			},
		},
		"code collision then success": {
			url: rawURL,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Twice()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					Return(model.Link{}, repolink.ErrConflict).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, l model.Link) (model.Link, error) { return l, nil }).Once()
			},
		},
		"conflict is url race returns existing": {
			url: rawURL,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					Return(model.Link{}, repolink.ErrConflict).Once()
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{Code: "RACED12"}, nil).Once()
			},
			wantCode: "RACED12",
		},
		"retries exhausted": {
			url: rawURL,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound)
				m.EXPECT().Create(mock.Anything, mock.Anything).
					Return(model.Link{}, repolink.ErrConflict)
			},
			wantErrStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			repo := repolink.NewMockRepository(t)
			tc.setup(repo)

			link, err := New(repo).Encode(context.Background(), tc.url)

			if tc.wantErrStatus != 0 {
				var ae *apperror.Error
				require.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantErrStatus, ae.HTTPStatus())
				return
			}
			require.NoError(t, err)
			if tc.wantCode != "" {
				assert.Equal(t, tc.wantCode, link.Code)
			} else {
				assert.Len(t, link.Code, model.ShortCodeLength)
			}
		})
	}
}
