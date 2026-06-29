package link

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

func TestDecode(t *testing.T) {
	errBoom := errors.New("boom")

	tcs := map[string]struct {
		code          string
		setup         func(*repolink.MockRepository)
		wantURL       string
		wantErrStatus int   // for apperror results
		wantErrIs     error // for propagated errors
	}{
		"returns original url": {
			code: "abc1234",
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByCode(mock.Anything, "abc1234").
					Return(model.Link{Code: "abc1234", OriginalURL: "https://example.com"}, nil).Once()
			},
			wantURL: "https://example.com",
		},
		"not found maps to 404": {
			code: "missing",
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByCode(mock.Anything, "missing").
					Return(model.Link{}, repolink.ErrNotFound).Once()
			},
			wantErrStatus: http.StatusNotFound,
		},
		"other repo error is propagated": {
			code: "abc1234",
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByCode(mock.Anything, "abc1234").
					Return(model.Link{}, errBoom).Once()
			},
			wantErrIs: errBoom,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			repo := repolink.NewMockRepository(t)
			tc.setup(repo)

			link, err := New(repo).Decode(context.Background(), tc.code)

			if tc.wantErrStatus != 0 {
				var ae *apperror.Error
				require.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantErrStatus, ae.HTTPStatus())
				return
			}
			if tc.wantErrIs != nil {
				require.ErrorIs(t, err, tc.wantErrIs)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantURL, link.OriginalURL)
		})
	}
}
