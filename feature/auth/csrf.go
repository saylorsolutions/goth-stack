package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/saylorsolutions/x/httpx"
	"net/http"
	"time"
	"yourapp/feature/audit"
)

const (
	CSRFFormKey   = "csrf_token"
	CSRFCookieKey = "_csrf"
)

type csrfKeyType string

var csrfKey csrfKeyType = CSRFCookieKey

func GetCSRF(r *http.Request) (string, bool) {
	value := r.Context().Value(csrfKey)
	if value == nil {
		return "", false
	}
	if sValue, ok := value.(string); ok {
		return sValue, true
	}
	return "", false
}

func (s *Service) SetCSRF() httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := make([]byte, 16)
			_, err := rand.Read(key)
			if err != nil {
				s.log.Postf(r.Context(), audit.AnonymousUser, "failed to read random key: %v", err)
				http.Error(w, "Failed to set CSRF token", 500)
				return
			}
			var csrfHex = make([]byte, len(key)*2)
			hex.Encode(csrfHex, key)
			if err := s.SetSecureCookie(w, CSRFCookieKey, string(csrfHex), 5*time.Minute); err != nil {
				http.Error(w, "Failed to set CSRF cookie", 500)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), csrfKey, string(csrfHex)))
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Service) RequireCSRF() httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := audit.AnonymousUser
			details, ok := GetSessionUser(r)
			if ok {
				username = details.Username
			}
			csrfValue, err := s.GetCookieValue(r, string(csrfKey))
			if err != nil {
				s.log.Post(r.Context(), username, "missing CSRF token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			givenCSRF := r.FormValue(CSRFFormKey)
			if csrfValue != givenCSRF {
				s.log.Post(r.Context(), username, "mismatched CSRF token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			s.log.Post(r.Context(), username, "Valid CSRF token")
			s.ClearCookie(w, string(csrfKey))
			next.ServeHTTP(w, r)
		})
	}
}
