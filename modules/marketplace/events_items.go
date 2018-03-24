package marketplace

import (
	"fmt"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func EventNewItem(item Item) {
	var (
		user      = item.User
		marketUrl = MARKETPLACE_SETTINGS.Sitename
		text      = fmt.Sprintf("[@%s](%s/user/%s) has added new item in category *%s* [%s](%s/user/%s/item/%s)",
			user.Username, marketUrl, user.Username,
			item.ItemCategory.NameEn,
			item.Name, marketUrl, user.Username, item.Uuid,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookItems, text)
}
