package marketplace

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) ViewShoutboxGET(w web.ResponseWriter, r *web.Request) {

	c.SelectedSection = "shoutbox"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	switch c.SelectedSection {
	case "shoutbox":
		thread, err := GetShoutboxThread(c.ViewUser.Language)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		if len(c.ViewThread.Messages) > 30 {
			c.ViewThread.Messages = c.ViewThread.Messages[0:30]
		}
	case "news":
		thread, err := GetNewsThread(c.ViewUser.Language)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	}

	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "shoutbox", c)
}

func (c *Context) ViewShoutboxPOST(w web.ResponseWriter, r *web.Request) {
	thread, err := GetShoutboxThread(c.ViewUser.Language)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewShoutboxGET(w, r)
		return
	}

	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewShoutboxGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewShoutboxGET(w, r)
		return
	}

	EventNewShoutboxPost(*c.ViewUser.User, *message)
	c.ViewShoutboxGET(w, r)
}
