package marketplace

import (
	"fmt"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func EventNewTransaction(transaction Transaction) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has created a new transaction [%s](%s/payments/%s)",
			transaction.Buyer.Username, marketUrl, transaction.Buyer.Username,
			transaction.Description+" | "+transaction.Type+" | "+transaction.Uuid, marketUrl, transaction.Uuid,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTransactions, text)
}

func EventTransactionStatusChange(transaction Transaction) {
	if transaction.CurrentPaymentStatus() == "PENDING" {
		return
	}
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("Transaction [%s](%s/payments/%s) (amount: %f / ) has changed status to **%s**",
			transaction.Description+" - "+transaction.Type+" "+transaction.Uuid, marketUrl, transaction.Uuid,
			transaction.CurrentAmountPaid(),
			transaction.CurrentPaymentStatus(),
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTransactions, text)
}
