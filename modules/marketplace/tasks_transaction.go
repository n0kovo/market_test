package marketplace

import (
	"time"

	"github.com/helloyi/go-waitgroup"
	"github.com/n0kovo/market_test/modules/util"
)

// TaskCleanInactiveReservations is a cron job that runs every 5 minutes.
// Deletes inactive reservations.
func TaskCleanInactiveReservations() {
	for range time.Tick(1 * time.Minute) {
		inactiveReservations := FindInactiveReservations()
		util.Log.Debug("[Task] [TaskCleanInactiveReservations] # of items: %d", len(inactiveReservations))
		for _, r := range inactiveReservations {
			transaction := r.Transaction
			if transaction.IsFailed() && transaction.IsCancelled() {
				transaction.SetTransactionStatus(
					"FAILED",
					transaction.CurrentAmountPaid(),
					"Transaction failed because no coins were transferred",
					"",
					nil,
				)
			}
			r.Remove()
		}
	}
}

// TaskUpdatePendingTransactions is a cron job that runs every 5 minutes.
// Checks for balance and updates transaction status.
func TaskUpdatePendingTransactions() {
	exec := func() {
		it := FindPendingTransactions()
		util.Log.Debug("[Task] [TaskUpdatePendingTransactions] # of items: %d", len(it))
		for _, t := range it {
			t.UpdateTransactionStatus()
		}
	}

	exec()
	for range time.Tick(1 * time.Minute) {
		exec()
	}
}

// TaskFailOldPendingTransactions is a cron job that runs every 5 minutes.
// Updates status of old pending payments
func TaskFailOldPendingTransactions() {
	exec := func() {
		it := FindOldPendingTransactions()
		util.Log.Debug("[Task] [TaskFailOldPendingTransactions] # of items: %d", len(it))
		for _, t := range it {
			t.Fail("Escrow failed automatically", "")
		}
	}

	exec()
	for range time.Tick(5 * time.Minute) {
		exec()
	}
}

// TaskReleaseConfirmedTransactions is a cron job that runs every 5 minutes.
// Releases funds to seller of confirmed transaction.
func TaskReleaseOldCompletedTransactions() {
	exec := func() {
		it := FindOldCompletedTransactions()
		util.Log.Debug("[Task] [TaskReleaseOldCompletedTransactions] # of items: %d", len(it))
		for _, t := range it {
			t.Release("Escrow released automatically", "")
		}
	}

	for range time.Tick(5 * time.Minute) {
		exec()
	}
}

// TaskDeleteFailedTransactions is a cron job that runs every 1 minute.
// Deletes failed transactions older than 10 days.
func TaskDeleteOldFailedTransactions() {
	for range time.Tick(1 * time.Minute) {
		it := FindOldFailedTransactions()
		util.Log.Debug("[Task] [TaskDeleteFailedTransactions] # of items: %d", len(it))
		for _, t := range it {
			t.Remove()
		}
	}
}

func TaskDeleteOldReleasedTransactions() {
	for range time.Tick(1 * time.Minute) {
		it := FindOldReleasedTransactions()
		util.Log.Debug("[Task] [TaskDeleteReleasedTransactions] # of items: %d", len(it))
		for _, t := range it {
			t.Remove()
		}
	}
}

func TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions() {
	wg := waitgroup.Create(12)
	ts := FindRecentlyCancelledAndReleasedTransactions()
	util.Log.Debug("[Task] [TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions] # of items: %d", len(ts))

	for i, _ := range ts {
		t := ts[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			defer wg.Done()
			t.UpdateTransactionStatus()
		}(&t)
	}

	wg.Wait()

}

func TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactionsCron() {
	TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions()
	for range time.Tick(60 * time.Minute) {
		TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions()
	}
}

func TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount() {

	it := FindReleasedAndCancelledTransactionsWithNonZeroAmount()
	wg := waitgroup.Create(1)

	util.Log.Debug("[Task] [TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount] # of items: %d", len(it))

	for i, _ := range it {
		t := it[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			defer wg.Done()
			if t.CurrentPaymentStatus() == "CANCELLED" {
				t.Cancel("Tx cancelled", "")
			}
			if t.CurrentPaymentStatus() == "RELEASED" {
				t.Release("Tx released", "")
			}
			t.UpdateTransactionStatus()
		}(&t)
	}

	wg.Wait()
}

func TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmountCron() {
	TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount()
	for range time.Tick(10 * time.Minute) {
		TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount()
	}
}

func TaskMaintainTransactions() {
	go TaskUpdatePendingTransactions()
	go TaskReleaseOldCompletedTransactions()
	go TaskFailOldPendingTransactions()
	go TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactionsCron()
	go TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmountCron()
}
