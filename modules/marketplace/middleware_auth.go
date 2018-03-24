package marketplace

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gocraft/web"
)

var (
	authorizedUrls = map[string]bool{
		"/":                  true,
		"/auth/login":        true,
		"/auth/recover":      true,
		"/api/auth/login":    true,
		"/api/auth/register": true,
		"/favicon.ico":       true,
		"/bot-check":         true,
	}
	passthruUrls = map[string]bool{
		"/item-image": true,
	}
	botCheckUuids = map[string]bool{}
)

func (c *Context) AuthMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	var userUuid string

	if len(r.URL.Query()["token"]) > 0 {
		apiSession, err := FindAPISessionByToken(r.URL.Query()["token"][0])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.APISession = apiSession
		if !c.APISession.IsTwoFactorSession || (c.APISession.IsTwoFactorSession && c.APISession.IsSecondFactorCompleted) {
			userUuid = apiSession.UserUuid
		}

	} else {
		session, _ := sessionStore.Get(r.Request, "auth-session")
		if session.Values["UserUuid"] != nil {
			userUuid = session.Values["UserUuid"].(string)
		}
	}

	c.ViewUser = User{}.ViewUser(c.Language)
	if passthruUrls[r.URL.Path] ||
		strings.HasPrefix(r.URL.Path, "/item-image") ||
		strings.HasPrefix(r.URL.Path, "/user-avatar") {
		next(w, r)
		return
	}

	if userUuid != "" {
		user, _ := FindUserByUuid(userUuid, false)
		if user == nil || user.Banned {
			http.NotFound(w, r.Request)
			return
		}

		var (
			oneMinute, _ = time.ParseDuration("1m")
			now          = time.Now()
		)
		if user.LastLoginDate.Add(oneMinute).Before(now) {
			user.LastLoginDate = &now
			user.Save()
		}

		c.ViewUser = user.ViewUser(user.Language)
	} else if !authorizedUrls[r.URL.Path] &&
		!strings.HasPrefix(r.URL.Path, "/captcha") &&
		!strings.HasPrefix(r.URL.Path, "/help") &&
		!strings.HasPrefix(r.URL.Path, "/item-image") &&
		!strings.HasPrefix(r.URL.Path, "/auth/register") {
		http.Redirect(w, r.Request, "/auth/login", 302)
		return
	}

	if c.ViewUser.User == nil {
		c.ViewUser.User = &User{}
	}

	next(w, r)
}

func (c *Context) BotCheckMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	if len(r.URL.Query()["json"]) > 0 {
		next(w, r)
		return
	}

	session, _ := sessionStore.Get(r.Request, "auth-session")
	checkUuid := session.Values["BotCheckUuid"]

	if checkUuid != nil {
		if _, ok := botCheckUuids[checkUuid.(string)]; ok {
			next(w, r)
		} else if !strings.HasPrefix(r.URL.Path, "/bot-check") && !strings.HasPrefix(r.URL.Path, "/captcha") {
			http.Redirect(w, r.Request, "/bot-check?redirect="+template.URLQueryEscaper(r.URL.String()), 302)
		} else {
			next(w, r)
		}
	} else if !strings.HasPrefix(r.URL.Path, "/bot-check") && !strings.HasPrefix(r.URL.Path, "/captcha") {
		http.Redirect(w, r.Request, "/bot-check?redirect="+template.URLQueryEscaper(r.URL.String()), 302)
	} else {
		next(w, r)
	}
}

func (c *Context) AdminMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}
	next(w, r)
}

func (c *Context) StaffMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.IsStaff || c.ViewUser.IsAdmin {
		next(w, r)
		return
	}
	http.NotFound(w, r.Request)
}
