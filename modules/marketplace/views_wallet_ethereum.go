package marketplace

import (
	"net/http"
	"strconv"

	btcqr "github.com/n0kovo/go.btcqr"
	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) EthereumWalletRecieve(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "wallet/ethereum/recieve", c)
}

func (c *Context) EthereumWalletSendGET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "wallet/ethereum/send", c)
}

func (c *Context) EthereumWalletSendPOST(w web.ResponseWriter, r *web.Request) {

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.EthereumWalletSendGET(w, r)
		return
	}

	var (
		address   = r.FormValue("address")
		amountStr = r.FormValue("amount")
	)

	if !ethereumRegexp.MatchString(address) {
		c.Error = "Wrong ETH address"
		c.EthereumWalletSendGET(w, r)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.Error = "Wrong amount"
		c.EthereumWalletSendGET(w, r)
		return
	}

	results, err := c.UserEthereumWallet.Send(address, amount)
	if err != nil {
		c.Error = err.Error()
		c.EthereumWalletSendGET(w, r)
		return
	}

	c.ETHPaymentResults, err = results.ETHPaymentResults()
	util.RenderTemplate(w, "wallet/ethereum/send_receipt", c)
}

func (c *Context) EthereumWalletImage(w web.ResponseWriter, r *web.Request) {
	req := &btcqr.Request{
		Address: r.PathParams["address"],
	}
	code, err := req.GenerateQR()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	png := code.PNG()
	w.Header().Set("Content-type", "image/png")
	w.Write(png)
}

func (c *Context) EthereumWalletActions(w web.ResponseWriter, r *web.Request) {
	c.UserEthereumWalletActions = FindUserEthereumWalletActionsForUser(c.ViewUser.Uuid)
	util.RenderTemplate(w, "wallet/ethereum/actions", c)
}
