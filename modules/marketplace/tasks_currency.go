package marketplace

import (
	"time"
)

func TaskUpdateCurrencyRates() {
	UpdateCurrencyRates()
	c := time.Tick(60 * time.Second)
	for range c {
		UpdateCurrencyRates()
	}
}
