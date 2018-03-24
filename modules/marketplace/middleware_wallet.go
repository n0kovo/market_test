package marketplace

import (
	"github.com/gocraft/web"
)

func (c *Context) BitcoinWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	for _, w := range c.UserBitcoinWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserBitcoinWallets) > 0 {
		c.UserBitcoinWallet = c.UserBitcoinWallets[0]
	}
	next(w, r)
}

func (c *Context) BitcoinCashWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	for _, w := range c.UserBitcoinCashWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserBitcoinCashWallets) > 0 {
		c.UserBitcoinCashWallet = c.UserBitcoinCashWallets[0]
	}
	next(w, r)
}

func (c *Context) EthereumWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	for _, w := range c.UserEthereumWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserEthereumWallets) > 0 {
		c.UserEthereumWallet = c.UserEthereumWallets[0]
	}
	next(w, r)
}

func (c *Context) WalletsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.Uuid != "" {
		c.UserBitcoinWallets = c.ViewUser.FindUserBitcoinWallets()
		c.UserEthereumWallets = c.ViewUser.FindUserEthereumWallets()
		c.UserBitcoinCashWallets = c.ViewUser.FindUserBitcoinCashWallets()
		c.UserBitcoinBalance = c.UserBitcoinWallets.Balance()
		c.UserEthereumBalance.Balance = c.UserEthereumWallets.Balance().Balance
		c.UserBitcoinCashBalance = c.UserBitcoinCashWallets.Balance()
	}

	next(w, r)
}
