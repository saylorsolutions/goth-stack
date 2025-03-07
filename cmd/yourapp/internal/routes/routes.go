package routes

import (
	"database/sql"
	"embed"
	"github.com/saylorsolutions/x/httpx"
	"log"
	"net/http"
	"yourapp/cmd/yourapp/internal/templates"
	"yourapp/feature/auth"
	"yourapp/feature/model"
	"yourapp/foundation/urlprefix"
)

type Router struct {
	Log          *log.Logger
	AuthSvc      *auth.Service
	Pool         *sql.DB
	StaticAssets embed.FS
	mux          *http.ServeMux
}

func (ro *Router) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	requireSession := ro.AuthSvc.RequireSession()
	requireAdmin := ro.requireAdmin()
	setCSRF := ro.AuthSvc.SetCSRF()
	requireCSRF := ro.AuthSvc.RequireCSRF()
	staticHandler := httpx.EmbeddedHandler(ro.StaticAssets, "", "")
	mux.HandleFunc("/blank", func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("GET /static/", staticHandler)
	mux.Handle("GET /{$}", requireSession(ro.homePage()))
	mux.Handle("GET /unauthorized", requireSession(ro.unauthorized()))
	mux.Handle("GET /pool", requireAdmin(ro.poolStats()))
	mux.Handle("GET /login", setCSRF(ro.loginPage()))
	mux.Handle("POST /login", requireCSRF(ro.loginHandling()))
	mux.Handle("/logout", requireSession(ro.logoutHandling()))
	return mux
}

func (ro *Router) homePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, ok := ro.getDetailsOrRedirect(w, r)
		if !ok {
			return
		}
		ro.Log.Println("Found auth details in request, username:", details.Username)
		ro.renderComponent(w, r, templates.Frame("", details.Username))
	}
}

func (ro *Router) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfVal, ok := auth.GetCSRF(r)
		if !ok {
			http.Error(w, "Missing CSRF token", 500)
			return
		}
		ro.renderComponent(w, r, templates.LoginPage(csrfVal))
	}
}

func (ro *Router) loginHandling() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		result, err := model.CheckPassword(r.Context(), ro.Pool, username, password)
		if err != nil {
			ro.Log.Println("Failed to check password:", err)
			w.WriteHeader(500)
			return
		}
		if !result.Matches {
			ro.AuthSvc.AuditEvent(r.Context(), username, "failed authentication")
			ro.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if _, err = ro.AuthSvc.SetAuthenticatedSession(w, r, username); err != nil {
			ro.Log.Println("Failed to set authenticated auth:", err)
			w.WriteHeader(500)
			return
		}
		ro.AuthSvc.AuditEvent(r.Context(), username, "Authenticated")
		ro.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func (ro *Router) logoutHandling() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, ok := auth.GetSessionUser(r)
		if !ok {
			ro.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ro.AuthSvc.ClearCookie(w, auth.SessionCookieName)
		if _, err := model.InvalidateSession(r.Context(), ro.Pool, details.SessionKey); err != nil {
			ro.Log.Println("[ERR] Failed to invalidate auth:", err)
		}
		ro.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (ro *Router) poolStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, ok := ro.getDetailsOrRedirect(w, r)
		if !ok {
			return
		}
		stats := ro.Pool.Stats()
		ro.renderComponent(w, r, templates.StatsPage(details.Username, stats))
	}
}

func (ro *Router) unauthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ro.renderComponent(w, r, templates.UnauthorizedPage())
	}
}

func (ro *Router) Redirect(w http.ResponseWriter, r *http.Request, location string, status int) {
	http.Redirect(w, r, urlprefix.Apply(location), status)
}
