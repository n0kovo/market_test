package marketplace

import (
	"fmt"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventUserLoggedIn(user User) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has logged in", user.Username, marketUrl, user.Username)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookAuthentication, text)
}

func EventUserRegistred(user User) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has registered", user.Username, marketUrl, user.Username)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookAuthentication, text)
}
