package marketplace

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"github.com/n0kovo/market_test/modules/apis"
	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) EditAdvertisings(w web.ResponseWriter, r *web.Request) {
	ads, err := FindAdvertisingByVendor(c.ViewUser.Uuid)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.Advertisings = ads
	c.AdvertisingCost = MARKETPLACE_SETTINGS.AdvertisingCost
	c.Items = FindItemsForSeller(c.ViewUser.Uuid)
	c.ViewSeller = Seller{c.ViewUser.User}.ViewSeller(c.ViewUser.User.Language)
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")
	util.RenderTemplate(w, "advertising/edit", c)
}

func (c *Context) AddAdvertisingsPOST(w web.ResponseWriter, r *web.Request) {
	count := 1000
	vendorUuid := c.ViewUser.Uuid
	comment := r.FormValue("text")
	itemUuid := r.FormValue("item")

	priceUSD := MARKETPLACE_SETTINGS.AdvertisingCost
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")

	price := priceUSD / c.USDBTCRate

	userWallets := c.ViewUser.User.FindUserBitcoinWallets()
	if userWallets.Balance().Balance < price {
		c.Error = fmt.Sprintf("Please deposit %f BTC to your onsite wallet.", price)
		c.EditAdvertisings(w, r)
		return
	}

	addr, err := apis.GenerateBTCAddress("advertising")
	if err != nil {
		c.Error = err.Error()
		c.EditAdvertisings(w, r)
		return
	}

	_, err = userWallets.Send(addr, price)
	if err != nil {
		c.Error = err.Error()
		c.EditAdvertisings(w, r)
		return
	}

	err = addAdvertising(comment, count, vendorUuid, itemUuid)
	if err != nil {
		c.Error = err.Error()
		//c.EditItem(w, r)
		c.EditAdvertisings(w, r)
		return
	}

	http.Redirect(w, r.Request, "/seller/"+c.ViewSeller.Username+"/advertisings", 302)

}
