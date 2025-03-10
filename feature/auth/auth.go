// Package auth provides http handler wrappers for repeatable logic.
package auth

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/saylorsolutions/x/env"
	"github.com/saylorsolutions/x/httpx"
	"net/http"
	"strings"
	"yourapp/feature/audit"
	"yourapp/feature/model"
)

const (
	expectedHashLen = 32
)

type Details struct {
	UserID     uint64
	Username   string
	Admin      bool
	SessionKey string
	Authz      []*model.UserAuthResult
}

func (d Details) HasAuth(auth string) bool {
	auth = strings.ToLower(auth)
	for _, granted := range d.Authz {
		if strings.ToLower(granted.Auth) == auth {
			return true
		}
	}
	return false
}

type Service struct {
	log      *audit.Logger
	sc       *securecookie.SecureCookie
	pool     *sql.DB
	userRepo model.UsersRepo
}

func initSecureCookie() (*securecookie.SecureCookie, error) {
	var hashKey []byte
	envHashKey := env.Val("SESSION_HASHKEY", "")
	if len(envHashKey) == 0 {
		hashKey = securecookie.GenerateRandomKey(expectedHashLen)
	} else {
		hashKey = make([]byte, expectedHashLen)
		size, err := hex.Decode(hashKey, []byte(envHashKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment provided hash key as hex: %w", err)
		}
		if size < expectedHashLen {
			return nil, fmt.Errorf("hash key must be a hex string representing at least %d bytes", expectedHashLen)
		}
	}
	return securecookie.New(hashKey, nil), nil
}

func NewAuthService(log *audit.Logger, pool *sql.DB) (*Service, error) {
	sc, err := initSecureCookie()
	if err != nil {
		return nil, err
	}
	return &Service{log: log, sc: sc, pool: pool}, nil
}

func (s *Service) RequireAuth(auth string) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			details, ok := GetSessionUser(r)
			if !ok {
				s.log.Postf(r.Context(), audit.AnonymousUser, "User is not granted auth '%s'", auth)
				http.Redirect(w, r, NoSessionRedirect, http.StatusFound)
				return
			}
			if !details.HasAuth(auth) {
				s.log.Postf(r.Context(), details.Username, "User is not granted auth '%s'", auth)
				http.Redirect(w, r, NoSessionRedirect, http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Service) AuditEvent(ctx context.Context, username string, message string) {
	s.log.Post(ctx, username, message)
}

func (s *Service) AuditEventf(ctx context.Context, username string, message string, args ...any) {
	s.log.Postf(ctx, username, message, args...)
}
