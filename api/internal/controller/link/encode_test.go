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

func TestEncode(t *testing.T) {
	const (
		rawURL     = "https://Example.com/Path"
		normalized = "https://example.com/Path"
		ip         = "203.0.113.7"
	)
	userID := uuid.New()

	echoCreate := func(_ context.Context, l model.Link) (model.Link, error) { return l, nil }

	tcs := map[string]struct {
		userID        *uuid.UUID
		setup         func(*repolink.MockRepository)
		wantCode      string // "" = random code, only length is checked
		wantErrStatus int    // 0 = no error
		wantOwned     bool   // result carries the user id
		wantCreatorIP bool   // result carries the client ip
	}{
		"dedup returns existing without touching quota": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{Code: "OLDCODE"}, nil).Once()
			},
			wantCode: "OLDCODE",
		},
		"anonymous under quota creates with creator ip": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().CountByCreatorIPSince(mock.Anything, ip, mock.Anything).Return(3, nil).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).RunAndReturn(echoCreate).Once()
			},
			wantCreatorIP: true,
		},
		"authenticated user skips quota and owns the link": {
			userID: &userID,
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).RunAndReturn(echoCreate).Once()
			},
			wantOwned: true,
		},
		"anonymous over quota returns 429": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().CountByCreatorIPSince(mock.Anything, ip, mock.Anything).
					Return(model.MaxAnonLinksPerDay, nil).Once()
			},
			wantErrStatus: http.StatusTooManyRequests,
		},
		"code collision then success": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Twice()
				m.EXPECT().CountByCreatorIPSince(mock.Anything, ip, mock.Anything).Return(0, nil).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					Return(model.Link{}, repolink.ErrConflict).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).RunAndReturn(echoCreate).Once()
			},
			wantCreatorIP: true,
		},
		"conflict is url race returns existing": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound).Once()
				m.EXPECT().CountByCreatorIPSince(mock.Anything, ip, mock.Anything).Return(0, nil).Once()
				m.EXPECT().Create(mock.Anything, mock.Anything).
					Return(model.Link{}, repolink.ErrConflict).Once()
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{Code: "RACED12"}, nil).Once()
			},
			wantCode: "RACED12",
		},
		"retries exhausted": {
			setup: func(m *repolink.MockRepository) {
				m.EXPECT().GetByNormalizedURL(mock.Anything, normalized).
					Return(model.Link{}, repolink.ErrNotFound)
				m.EXPECT().CountByCreatorIPSince(mock.Anything, ip, mock.Anything).Return(0, nil).Once()
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

			link, err := New(repo).Encode(context.Background(), rawURL, tc.userID, ip)

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
			if tc.wantOwned {
				require.NotNil(t, link.UserID)
				assert.Equal(t, userID, *link.UserID)
				assert.Nil(t, link.CreatorIP)
			}
			if tc.wantCreatorIP {
				require.NotNil(t, link.CreatorIP)
				assert.Equal(t, ip, *link.CreatorIP)
				assert.Nil(t, link.UserID)
			}
		})
	}
}
