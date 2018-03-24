package marketplace

import (
	"github.com/gocraft/web"
)

func (c *Context) LocalizationMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	lang := c.ViewUser.Language
	if lang == "" {
		session, _ := sessionStore.Get(r.Request, "lang-session")
		if len(r.URL.Query()["language"]) > 0 {
			lang = r.URL.Query()["language"][0]

			session.Values["language"] = lang
			session.Save(r.Request, w)
		} else {
			langSess := session.Values["language"]
			if langSess != nil {
				lang = langSess.(string)
			}
		}
		c.Language = lang
		c.ViewUser.Language = lang
	}

	c.Localization = GetLocalization(lang)
	next(w, r)
}
