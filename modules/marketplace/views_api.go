package marketplace

import (
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewAPILogin(user User, w web.ResponseWriter, r *web.Request) {
	var err error
	c.APISession, err = CreateAPISession(user)
	if err != nil {
		c.Error = err.Error()
		util.APIResponse(w, r, c)
		return
	}

	now := time.Now()
	user.LastLoginDate = &now
	user.Save()

	EventUserLoggedIn(user)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPILoginRegisterGET(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}
	c.CaptchaId = captcha.New()
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPILoginPOST(w web.ResponseWriter, r *web.Request) {
	if r.FormValue("decryptedmessage") == "" {
		var (
			isCaptchaValid    = captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
			user, _           = FindUserByUsername(r.FormValue("username"))
			isLoginSuccessful = isCaptchaValid && (user != nil) && user.CheckPassphrase(r.FormValue("passphrase"))
		)
		if !isCaptchaValid {
			c.Error = "Invalid captcha"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
		if user == nil || !isLoginSuccessful {
			c.Error = "Failed to authenticate"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
		if user.TwoFactorAuthentication {
			session, _ := CreateAPISession(*user)
			c.APISession = session
			c.APISession.SecondFactorSecretText = util.GenerateUuid()
			c.APISession.Save()

			c.SecretText, _ = util.EncryptText(c.APISession.SecondFactorSecretText, user.Pgp)
			util.APIResponse(w, r, c)
		} else {
			c.ViewAPILogin(*user, w, r)
		}
	} else {
		var (
			secretText       = c.APISession.SecondFactorSecretText
			decryptedmessage = strings.Trim(r.FormValue("decryptedmessage"), "\n ")
		)
		if decryptedmessage == secretText {
			c.ViewAPILogin(c.APISession.User, w, r)
			return
		} else {
			c.Error = "Could not authenticate"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
	}
}

func (c *Context) ViewAPISERP(w web.ResponseWriter, r *web.Request) {
	c.listAvailableItems(w, r)
	util.APIResponse(w, r, c)
}
