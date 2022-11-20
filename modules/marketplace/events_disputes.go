package marketplace

import (
	"fmt"
	"strings"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventNewDispute(dispute Dispute) {
	transaction, _ := dispute.Transaction()
	if transaction == nil {
		return
	}
	if len(dispute.DisputeClaims) == 0 {
		return
	}
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		claim     = dispute.DisputeClaims[0]
		text      = fmt.Sprintf("[@%s](%s/user/%s) has created a new dispute [%s](%s/payments/%s)\n**Claim:**\n> %s\n\n**Status:**\n%s",
			transaction.Buyer.Username, marketUrl, transaction.Uuid,
			transaction.Buyer.Username, marketUrl, dispute.Uuid,
			strings.Replace(claim.Claim, "\n", "\n > ", -1), claim.Status,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookDisputes, text)
}

func EventDisputeNewMessage(dispute Dispute, message Message) {
	transaction, _ := dispute.Transaction()
	if transaction == nil {
		return
	}
	if len(dispute.DisputeClaims) == 0 {
		return
	}
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		claim     = dispute.DisputeClaims[0]
		text      = fmt.Sprintf("[@%s](%s/user/%s) has commented on a dispute [%s](%s/payments/%s)\n**Claim:**\n> %s\n\n**Status:**\n%s\n**Comment:**\n> %s\n",
			message.SenderUser.Username, marketUrl, message.SenderUser.Username,
			transaction.Buyer.Username, marketUrl, dispute.Uuid,
			strings.Replace(claim.Claim, "\n", "\n > ", -1), claim.Status,
			strings.Replace(message.Text, "\n", "\n > ", -1),
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookDisputes, text)
}
