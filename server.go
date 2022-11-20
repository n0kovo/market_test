package main

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	"github.com/gorilla/context"

	"github.com/n0kovo/market_test/modules/marketplace"
	"github.com/n0kovo/market_test/modules/settings"
)

func RunServer() {

	var (
		settings = settings.GetSettings()
	)

	// Crons
	if !settings.Debug {
		go marketplace.TaskMaintainTransactions()
		go marketplace.TaskMaintainWallets()
	}

	go marketplace.TaskUpdateCurrencyRates()

	// Root router
	rootRouter := web.New(marketplace.Context{})
	rootRouter.Middleware(web.StaticMiddleware("public"))

	// Common middleware
	// if !settings.Debug {
	// 	rootRouter.Middleware((*marketplace.Context).BotCheckMiddleware)
	// }

	rootRouter.Middleware((*marketplace.Context).AuthMiddleware)
	rootRouter.Middleware((*marketplace.Context).ModeMarketplaceMiddleware)
	rootRouter.Middleware((*marketplace.Context).LoggerMiddleware)
	rootRouter.Middleware((*marketplace.Context).LocalizationMiddleware)
	rootRouter.Middleware((*marketplace.Context).CurrencyMiddleware)

	// Marketplace routes
	marketplace.ConfigureRouter(rootRouter.Subrouter(marketplace.Context{}, "/"))

	// Start the server
	address := fmt.Sprintf("%s:%s", settings.Host, settings.Port)
	println("Running server on " + address)
	http.ListenAndServe(address, context.ClearHandler(rootRouter))

}
