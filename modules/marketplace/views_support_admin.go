package marketplace

import (
	"github.com/gocraft/web"
	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) ListSupportThreads(w web.ResponseWriter, r *web.Request) {
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	if c.SelectedSection == "unanswered" {
		b := false
		c.SupportThreads = FindSupportThreads(&b)
	} else {
		c.SupportThreads = FindSupportThreads(nil)
	}

	util.RenderTemplate(w, "support/admin/threads", c)
}
