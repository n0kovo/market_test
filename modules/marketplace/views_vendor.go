package marketplace

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ShowStore(w web.ResponseWriter, r *web.Request) {

	if !c.ViewSeller.IsSeller {
		redirectUrl := "/user/" + c.ViewSeller.Username + "/about"
		http.Redirect(w, r.Request, redirectUrl, 302)
		return
	}

	items := FindItemsForSeller(c.ViewSeller.Uuid)
	c.ViewItems = items.ViewItems(c.ViewUser.Language)

	// rendering
	util.RenderTemplate(w, "store/list_items", c)
}

func (c *Context) AboutStore(w web.ResponseWriter, r *web.Request) {
	if c.ViewSeller.IsTrustedSeller {
		th, err := GetVendorVerificationThread(*c.ViewSeller.User, false)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.ViewThread = th.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	}
	util.RenderTemplate(w, "store/about", c)
}

func (c *Context) StoreBoard(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "store/board", c)
}

func (c *Context) StoreBoardPost(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.StoreBoard(w, r)
		return
	}
	thread, err := GetStoreThread(*c.ViewSeller.Seller)
	if err != nil {
		c.Error = err.Error()
		c.StoreBoard(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.StoreBoard(w, r)
		return
	}
	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.StoreBoard(w, r)
		return
	}
	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)

	EventNewVendorMessageboardPost(*c.ViewUser.User, *c.ViewSeller.Seller, *message)
	c.StoreBoard(w, r)
}

func (c *Context) StoreReviews(w web.ResponseWriter, r *web.Request) {
	if !c.ViewSeller.IsSeller {
		redirectUrl := "/user/" + c.ViewSeller.Username + "/about"
		http.Redirect(w, r.Request, redirectUrl, 302)
		return
	}

	util.RenderTemplate(w, "store/reviews", c)

}
