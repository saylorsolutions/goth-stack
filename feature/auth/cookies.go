package auth

import (
	"net/http"
	"time"
)

func (s *Service) SetSecureCookie(w http.ResponseWriter, key string, value string, cookieTTL time.Duration) error {
	val, err := s.sc.Encode(key, value)
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		Expires:  time.Now().Add(cookieTTL),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (s *Service) GetCookieValue(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	var val string
	if err := s.sc.Decode(key, cookie.Value, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (s *Service) ClearCookie(w http.ResponseWriter, key string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
}
