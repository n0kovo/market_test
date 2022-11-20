package marketplace

import (
	"fmt"
	"strings"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventNewShoutboxPost(user User, message Message) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has logged posted in shoutbox:\n> %s",
			user.Username, marketUrl, user.Username,
			strings.Replace(message.Text, "\n", "\n > ", -1), //------------|
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookShoutbox, text)
}
