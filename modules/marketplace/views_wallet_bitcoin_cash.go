package marketplace

import (
	_ "math"
	"net/http"
	"strconv"
	"strings"

	btcqr "github.com/GeertJohan/go.btcqr"
	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) BitcoinCashWalletRecieve(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "wallet/bitcoin_cash/recieve", c)
}

func (c *Context) BitcoinCashWalletSendStep1GET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "wallet/bitcoin_cash/send_step_1", c)
}

func (c *Context) BitcoinCashWalletSendStep2GET(w web.ResponseWriter, r *web.Request) {
	var err error
	var withdraw bool = (r.FormValue("withdraw") == "1")
	c.BCHFee, c.Description, err = c.UserBitcoinCashWallets.EstimateFee(c.Address, c.Amount)
	if err != nil {
		c.Error = err.Error()
		c.BitcoinCashWalletSendStep1GET(w, r)
		return
	}
	if withdraw {
		c.Amount = c.UserBitcoinCashBalance.Balance // util.RoundPlus(math.Max(c.UserBitcoinCashBalance.Balance-c.BCHFee, 0.0), 8)
	}
	if c.UserBitcoinCashBalance.Balance < c.BCHFee {
		c.Error = "Amount + Fee is greater than balance of your wallet"
	}
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "wallet/bitcoin_cash/send_step_2", c)
}

func (c *Context) BitcoinCashWalletSendPOST(w web.ResponseWriter, r *web.Request) {
	switch r.FormValue("step") {
	case "1":
		c.BitcoinCashWalletSendStep1POST(w, r)
		break
	case "2":
		c.BitcoinCashWalletSendStep2POST(w, r)
		break
	default:
		http.NotFound(w, r.Request)
		return
	}
}

func (c *Context) BitcoinCashWalletSendStep1POST(w web.ResponseWriter, r *web.Request) {
	c.Address = strings.Trim(r.FormValue("address"), " ")
	amountStr := r.FormValue("amount")

	if !bitcoinRegexp.MatchString(c.Address) {
		c.Error = "Wrong BCH address"
		c.BitcoinCashWalletSendStep1GET(w, r)
		return
	}

	var err error
	if r.FormValue("withdraw") != "1" {
		c.Amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			c.Error = "Wrong amount"
			c.BitcoinCashWalletSendStep1GET(w, r)
			return
		}
	}

	c.BitcoinCashWalletSendStep2GET(w, r)
}

func (c *Context) BitcoinCashWalletSendStep2POST(w web.ResponseWriter, r *web.Request) {

	var (
		amountStr      = r.FormValue("amount")
		isCaptchaValid = captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
		err            error
	)

	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.BitcoinCashWalletSendStep2GET(w, r)
		return
	}

	c.Address = strings.Trim(r.FormValue("address"), " ")
	if !bitcoinRegexp.MatchString(c.Address) {
		c.Error = "Wrong BCH address **"
		c.BitcoinCashWalletSendStep2GET(w, r)
		return
	}

	c.Amount, err = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.Error = "Wrong amount"
		c.BitcoinCashWalletSendStep2GET(w, r)
		return
	}

	receipt, err := c.UserBitcoinCashWallets.Send(c.Address, c.Amount)
	if err != nil {
		c.Error = err.Error()
		util.RenderTemplate(w, "wallet/bitcoin_cash/send_step_2", c)
		return
	}

	c.BCHPaymentResult, err = receipt.BCHPaymentResult()
	util.RenderTemplate(w, "wallet/bitcoin_cash/send_step_3", c)
}

func (c *Context) BitcoinCashWalletImage(w web.ResponseWriter, r *web.Request) {
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

func (c *Context) BitcoinCashWalletActions(w web.ResponseWriter, r *web.Request) {
	c.UserBitcoinCashWalletActions = FindUserBitcoinCashWalletActionsForUser(c.ViewUser.Uuid)
	util.RenderTemplate(w, "wallet/bitcoin_cash/actions", c)
}
