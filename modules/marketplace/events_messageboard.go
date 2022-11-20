package marketplace

import (
	"fmt"
	"strings"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventNewVendorMessageboardPost(user User, vendor Seller, message Message) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has logged posted in thread [@%s](%s/user/%s/board):\n> %s",
			user.Username, marketUrl, user.Username,
			vendor.Username, marketUrl, vendor.Username,
			strings.Replace(message.Text, "\n", "\n > ", -1),
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookMessageboard, text)
}
