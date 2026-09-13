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

func TestResolve(t *testing.T) {
	const code = "abc1234"

	tcs := map[string]struct {
		setup         func(*repolink.MockRepository)
		wantURL       string
		wantErrStatus int
	}{
		"found returns link and counts the visit": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().IncrementAndGetByCode(mock.Anything, code).
					Return(model.Link{Code: code, OriginalURL: "https://x.com", ClickCount: 5}, nil).Once()
			},
			wantURL: "https://x.com",
		},
		"not found maps to 404": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().IncrementAndGetByCode(mock.Anything, code).
					Return(model.Link{}, repolink.ErrNotFound).Once()
			},
			wantErrStatus: http.StatusNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			repo := repolink.NewMockRepository(t)
			tc.setup(repo)

			link, err := New(repo).Resolve(context.Background(), code)
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
