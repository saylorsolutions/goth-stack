package routes

import (
	"bytes"
	"github.com/a-h/templ"
	"github.com/saylorsolutions/x/httpx"
	"io"
	"net/http"
	"yourapp/feature/auth"
)

func (ro *Router) renderComponent(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	var buf bytes.Buffer
	if err := comp.Render(r.Context(), &buf); err != nil {
		ro.Log.Println("Failed to render component for path:", r.URL.Path)
		http.Error(w, err.Error(), 500)
	}
	_, _ = io.Copy(w, &buf)
}

func (ro *Router) getDetailsOrRedirect(w http.ResponseWriter, r *http.Request) (auth.Details, bool) {
	details, ok := auth.GetSessionUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return auth.Details{}, false
	}
	return details, true
}

func (ro *Router) requireAdmin() httpx.Middleware {
	return func(next http.Handler) http.Handler {
		session := ro.AuthSvc.RequireSession()
		return session(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			details, ok := ro.getDetailsOrRedirect(w, r)
			if !ok {
				return
			}
			if !details.Admin {
				ro.AuthSvc.AuditEventf(r.Context(), details.Username, "Attempted to access %s %s, not admin", r.Method, r.URL.Path)
				http.Redirect(w, r, "/unauthorized", http.StatusFound)
			}
			next.ServeHTTP(w, r)
		}))
	}
}

func (ro *Router) requireAuth(auth string) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			details, ok := ro.getDetailsOrRedirect(w, r)
			if !ok {
				return
			}
			if !details.HasAuth(auth) {
				ro.AuthSvc.AuditEventf(r.Context(), details.Username, "Attempted to access %s %s, not granted auth %s", r.Method, r.URL.Path, auth)
				http.Redirect(w, r, "/unauthorized", http.StatusFound)
			}
			next.ServeHTTP(w, r)
		})
	}
}
