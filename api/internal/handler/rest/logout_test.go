package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrlauth "github.com/kytruongdev/shortlink/internal/controller/auth"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

func TestLogoutHandler(t *testing.T) {
	ctrl := ctrlauth.NewMockController(t)
	ctrl.EXPECT().Logout(mock.Anything, "some-refresh").Return(nil).Once()

	h := NewAuth(ctrl, false, time.Hour)
	wrapped := httpserver.HandlerErr(h.Logout)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: refreshCookieName, Value: "some-refresh"})
	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	cookie := findCookie(rec.Result().Cookies(), refreshCookieName)
	require.NotNil(t, cookie)
	assert.True(t, cookie.MaxAge < 0, "logout expires the refresh cookie")
}
