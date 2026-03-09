package auth

import (
	"context"
	"net/http"
	"time"
	"yourapp/feature/audit"
	"yourapp/foundation/urlprefix"

	"github.com/saylorsolutions/x/httpx"
)

const (
	SessionCookieName = "JSESSIONID"
)

var (
	NoSessionRedirect = urlprefix.Apply("/login")
)

type sessionDetails string

const sessionDetailsKey = sessionDetails("sessionDetails")

func GetSessionUser(r *http.Request) (Details, bool) {
	val := r.Context().Value(sessionDetailsKey)
	if val == nil {
		return Details{}, false
	}
	details, ok := val.(Details)
	if !ok {
		return Details{}, false
	}
	return details, true
}

func setSessionDetails(r *http.Request, details Details) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), sessionDetailsKey, details))
}

func (s *Service) RequireSession() httpx.Middleware {
	mw := httpx.DeferMiddleware()
	return func(next http.Handler) http.Handler {
		return mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionKey, err := s.GetCookieValue(r, SessionCookieName)
			if err != nil {
				http.Redirect(w, r, NoSessionRedirect, http.StatusFound)
				return
			}
			result, err := s.userRepo.GetSessionUser(r.Context(), sessionKey)
			if err != nil || result == nil {
				s.log.Postf(r.Context(), audit.AnonymousUser, "Failed to get user details for auth %s: %v", sessionKey, err)
				http.Redirect(w, r, NoSessionRedirect, http.StatusFound)
				return
			}
			if len(result.Username) == 0 {
				s.log.Postf(r.Context(), audit.AnonymousUser, "Failed to get user details for auth: %s", sessionKey)
				http.Redirect(w, r, NoSessionRedirect, http.StatusFound)
				return
			}
			details := Details{
				UserID:     result.UserID,
				Username:   result.Username,
				Admin:      result.Admin,
				SessionKey: sessionKey,
			}
			_, err = s.userRepo.UpdateSessionLiveness(r.Context(), sessionKey)
			if err != nil {
				s.log.Postf(r.Context(), details.Username, "Failed to update session %s liveness: %v", sessionKey, err)
				http.Error(w, "Session management error", 500)
				return
			}
			if err := s.setSessionCookie(w, sessionKey); err != nil {
				s.log.Postf(r.Context(), details.Username, "Failed to encode session key as cookie value: %v", err)
				http.Error(w, "Session management error", 500)
				return
			}
			auths, err := s.userRepo.UserAuth(r.Context(), result.UserID)
			if err != nil {
				s.log.Postf(r.Context(), audit.AnonymousUser, "Failed to retrieve authorizations for user %s: %v", result.Username, err)
				http.Error(w, "Session management error", 500)
				return
			}
			details.Authz = auths
			r = setSessionDetails(r, details)
			s.log.Postf(r.Context(), result.Username, "%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		}))
	}
}

func (s *Service) SetAuthenticatedSession(w http.ResponseWriter, r *http.Request, username string) (*http.Request, error) {
	ses, err := s.userRepo.CreateSession(r.Context(), username)
	if err != nil {
		return r, err
	}
	details, err := s.userRepo.GetSessionUser(r.Context(), ses.SessionKey)
	if err != nil {
		return r, err
	}
	userDetails := Details{
		Username:   details.Username,
		Admin:      details.Admin,
		SessionKey: ses.SessionKey,
	}
	authz, err := s.userRepo.UserAuth(r.Context(), details.UserID)
	if err != nil {
		return r, err
	}
	userDetails.Authz = authz
	if err := s.setSessionCookie(w, ses.SessionKey); err != nil {
		s.log.Postf(r.Context(), username, "Unable to encode session key as cookie value: %v", err)
		return r, err
	}
	r = setSessionDetails(r, userDetails)
	return r, nil
}

func (s *Service) setSessionCookie(w http.ResponseWriter, sessionKey string) error {
	return s.SetSecureCookie(w, SessionCookieName, sessionKey, 30*time.Minute)
}

func (s *Service) InvalidateSession(r *http.Request) error {
	details, ok := GetSessionUser(r)
	if !ok {
		return nil
	}
	if _, err := s.userRepo.InvalidateSession(r.Context(), details.SessionKey); err != nil {
		return err
	}
	return nil
}
