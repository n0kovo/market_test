package marketplace

import (
	"time"

	"github.com/helloyi/go-waitgroup"
)

func UpdateOldFailedAndPendingTransactions() {

	wg := waitgroup.Create(12)
	ts := findFailedPendingTransactionsCreatedGteTimestamp(time.Now().AddDate(0, 0, -30))
	println("[Command] [UpdateOldFailedAndPendingTransactions] # of items: ")

	for i, _ := range ts {
		t := ts[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			println(i, "/", len(ts))
			defer wg.Done()
			println("Updating transaction ", t.Uuid)
			t.UpdateTransactionStatus()
		}(&t)
	}

	wg.Wait()
}

func ResendReleasedTransactions() {
	wg := waitgroup.Create(12)
	ts := findReleasedTransactionsCreatedGteTimestamp(time.Now().AddDate(-1, 0, 0))
	for i, _ := range ts {
		t := ts[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			println(i, "/", len(ts))
			defer wg.Done()
			println("Updating transaction ", t.Uuid)
			t.Release("Routine transaction update", "")
		}(&t)
	}

	wg.Wait()
}
