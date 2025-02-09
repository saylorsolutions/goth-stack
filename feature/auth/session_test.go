package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/saylorsolutions/x/httpx"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"yourapp/feature/audit"
	"yourapp/feature/model"
)

func TestService_RequireSession(t *testing.T) {
	var (
		loginRedirectCalls         int
		getSessionCalls            int
		getAuthCalls               int
		authenticatedCalls         int
		updateSessionLivenessCalls int
	)
	resetCounts := func() {
		loginRedirectCalls = 0
		getSessionCalls = 0
		getAuthCalls = 0
		authenticatedCalls = 0
		updateSessionLivenessCalls = 0
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	authSvc := testAuthService(t, ctx)
	authSvc.userRepo.RedirectGetSessionUser(func(_ context.Context, _ *sql.DB, sessionKey string) (*model.GetSessionUserResult, error) {
		getSessionCalls++
		if sessionKey != "abc" {
			return nil, sql.ErrNoRows
		}
		return &model.GetSessionUserResult{
			UserID:   1,
			Username: "Bob",
			Admin:    false,
		}, nil
	})
	authSvc.userRepo.RedirectUpdateSessionLiveness(func(_ context.Context, _ *sql.DB, sessionKey string) (sql.Result, error) {
		updateSessionLivenessCalls++
		return nil, nil
	})
	authSvc.userRepo.RedirectUserAuth(func(_ context.Context, _ *sql.DB, user uint64) ([]*model.UserAuthResult, error) {
		getAuthCalls++
		return nil, nil
	})
	requireSession := authSvc.RequireSession()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		loginRedirectCalls++
	})
	mux.HandleFunc("GET /unprotected", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
	mux.Handle("GET /protected", requireSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedCalls++
	})))

	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Run("No auth needed", func(t *testing.T) {
		resetCounts()
		_, status, err := httpx.GetRequest(fmt.Sprintf("%s/unprotected", srv.URL)).Send()
		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, 0, loginRedirectCalls)
		assert.Equal(t, 0, getSessionCalls)
		assert.Equal(t, 0, getAuthCalls)
		assert.Equal(t, 0, updateSessionLivenessCalls)
		assert.Equal(t, 0, authenticatedCalls)
	})

	t.Run("Missing session cookie", func(t *testing.T) {
		resetCounts()
		_, status, err := httpx.GetRequest(fmt.Sprintf("%s/protected", srv.URL)).Send()
		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, 1, loginRedirectCalls)
		assert.Equal(t, 0, getSessionCalls)
		assert.Equal(t, 0, getAuthCalls)
		assert.Equal(t, 0, updateSessionLivenessCalls)
		assert.Equal(t, 0, authenticatedCalls)
	})

	t.Run("Invalid session", func(t *testing.T) {
		resetCounts()
		_, status, err := httpx.GetRequest(fmt.Sprintf("%s/protected", srv.URL)).AddHeader("Set-Cookie", "abc").Send()
		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, 1, loginRedirectCalls)
		assert.Equal(t, 0, getSessionCalls)
		assert.Equal(t, 0, getAuthCalls)
		assert.Equal(t, 0, updateSessionLivenessCalls)
		assert.Equal(t, 0, authenticatedCalls)
	})

	t.Run("Old session", func(t *testing.T) {
		resetCounts()
		// This relies on the fact that the query to return information from the session store will implicitly filter out old sessions and return no rows.
		val, err := authSvc.sc.Encode(SessionCookieName, "old-session-key")
		assert.NoError(t, err)
		_, status, err := httpx.GetRequest(fmt.Sprintf("%s/protected", srv.URL)).AddHeader("Set-Cookie", "value").
			SetCookie(&http.Cookie{
				Name:     SessionCookieName,
				Value:    val,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			}).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, 1, loginRedirectCalls)
		assert.Equal(t, 1, getSessionCalls)
		assert.Equal(t, 0, getAuthCalls)
		assert.Equal(t, 0, updateSessionLivenessCalls)
		assert.Equal(t, 0, authenticatedCalls)
	})

	t.Run("Valid session", func(t *testing.T) {
		resetCounts()
		val, err := authSvc.sc.Encode(SessionCookieName, "abc")
		assert.NoError(t, err)
		_, status, err := httpx.GetRequest(fmt.Sprintf("%s/protected", srv.URL)).
			SetCookie(&http.Cookie{
				Name:     SessionCookieName,
				Value:    val,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			}).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, 0, loginRedirectCalls)
		assert.Equal(t, 1, getSessionCalls)
		assert.Equal(t, 1, getAuthCalls)
		assert.Equal(t, 1, updateSessionLivenessCalls)
		assert.Equal(t, 1, authenticatedCalls)
	})
}

func testAuthService(t *testing.T, ctx context.Context) *Service {
	auditLog := audit.NewLogger(nil, audit.StdDelegate(log.Default(), true))
	auditLog.UserRepo.RedirectInsertAuditLog(func(_ context.Context, _ *sql.DB, user string, msg string) (sql.Result, error) {
		t.Log("[Audit Log]", user, msg)
		return nil, nil
	})
	sc := securecookie.New([]byte("abc"), nil)
	return &Service{
		log: auditLog,
		sc:  sc,
	}
}
