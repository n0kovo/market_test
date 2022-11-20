package marketplace

import (
	"fmt"
	"strings"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventNewTrustedVendorRequest(vendor Seller) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has requested for a trusted vendor status",
			vendor.Username, marketUrl, vendor.Username,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTrustedVendors, text)
}

func EventNewTrustedVendorThreadPost(user User, vendor Seller, message Message) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has posted in vendor verification thread [@%s](%s/staff/vendors/%s):\n> %s",
			user.Username, marketUrl, user.Username,
			vendor.Username, marketUrl, vendor.Username,
			strings.Replace(message.Text, "\n", "\n > ", -1), //------------|
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTrustedVendors, text)
}
