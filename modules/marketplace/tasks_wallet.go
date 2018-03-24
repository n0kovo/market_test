package marketplace

import (
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func TaskUpdateRecentBitcoinWallets() {
	exec := func() {
		recentWallets := FindRecentBitcoinWallets()
		util.Log.Debug("[Task] [TaskUpdateRecentBitcoinWallets], number of wallets to update: %d", len(recentWallets))
		for _, w := range recentWallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(60 * time.Minute) {
		exec()
	}
}

func TaskUpdateAllBitcoinWallets() {
	exec := func() {
		util.Log.Debug("[Task] [TaskUpdateAllBitcoinCashWallets]")
		wallets := FindAllBitcoinCashWallets()
		for _, w := range wallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(24 * time.Hour) {
		exec()
	}
}

func TaskUpdateRecentBitcoinCashWallets() {
	exec := func() {
		recentWallets := FindRecentBitcoinCashWallets()
		util.Log.Debug("[Task] [TaskUpdateRecentBitcoinCashWallets], number of wallets to update: %d", len(recentWallets))
		for _, w := range recentWallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(60 * time.Minute) {
		exec()
	}
}

func TaskUpdateAllBitcoinCashWallets() {
	exec := func() {
		util.Log.Debug("[Task] [TaskUpdateAllBitcoinCashWallets]")
		wallets := FindAllBitcoinCashWallets()
		for _, w := range wallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(24 * time.Hour) {
		exec()
	}
}

func TaskUpdateRecentEthereumWallets() {
	exec := func() {
		recentWallets := FindRecentEthereumWallets()
		util.Log.Debug("[Task] [TaskUpdateRecentEthereumWallets], number of wallets to update: %d", len(recentWallets))
		for _, w := range recentWallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(60 * time.Minute) {
		exec()
	}
}

func TaskUpdateAllEthereumWallets() {
	exec := func() {
		util.Log.Debug("[Task] [TaskUpdateAllEthereumWallets]")
		wallets := FindAllEthereumWallets()
		for _, w := range wallets {
			w.UpdateBalance(false)
		}
	}
	for range time.Tick(24 * time.Hour) {
		exec()
	}
}

func TaskMaintainWallets() {
	go TaskUpdateRecentBitcoinWallets()
	go TaskUpdateAllBitcoinWallets()
	go TaskUpdateRecentBitcoinCashWallets()
	go TaskUpdateAllBitcoinCashWallets()
	go TaskUpdateRecentEthereumWallets()
	go TaskUpdateAllEthereumWallets()
}
