package marketplace

import (
	"github.com/gocraft/web"
	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) ShowFeed(w web.ResponseWriter, r *web.Request) {
	if len(r.URL.Query()["section"]) > 0 {
		section := r.URL.Query()["section"][0]
		c.SelectedSection = section
	} else {
		c.SelectedSection = ""
	}
	feedItems := CacheGetPublicFeedItems()
	c.ViewFeedItems = feedItems.ViewFeedItems(c.ViewUser.Language, c.SelectedSection)
	util.RenderTemplate(w, "feed", c)
}
