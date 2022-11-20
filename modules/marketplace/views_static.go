package marketplace

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"github.com/n0kovo/market_test/modules/util"
	"strings"
)

func loadStaticPages(dirname string) ([]StaticPage, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	pages := []StaticPage{}

	for _, fileInfo := range files {
		basename := fileInfo.Name()
		if basename[0] == '.' {
			continue
		}
		fullname := path.Join(dirname, basename)
		contents, err := ioutil.ReadFile(fullname)
		if err != nil {
			continue
		}
		pageBody := template.HTML(blackfriday.MarkdownCommon([]byte(contents)))
		titleParts := strings.Split(basename, "_")

		page := StaticPage{
			Title:    strings.Join(titleParts[1:len(titleParts)], " "),
			Body:     pageBody,
			Filename: basename,
		}

		pages = append(pages, page)

	}
	return pages, nil
}

func (c *Context) Help(w web.ResponseWriter, r *web.Request) {
	lang := c.ViewUser.Language
	if lang == "" {
		lang = "en"
	}

	if lang != "en" && lang != "ru" {
		lang = "en"
	}

	pages, err := loadStaticPages(fmt.Sprintf("./static/%s/", lang))
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.StaticPages = pages
	c.StaticPage = StaticPage{
		Title: "Support Hub",
		Body: template.HTML(`
<pre>
████████████▄░░░░░░░░░░░░░░░░░░░░░░░░░░░░░████
████████████▄▄▄░░░░░░░░░░░░░░░░░░░░░░░░░░░░███
███████▀█▀███████▄▄░░░░░░░░░░░░░░░░░░░░░░░░░▀█
███████▀░░░░▄█████▄░░░░░░░░░░░░░░░░░░░░░░░░░░█
████████▄░░░░▀██████░░░░░░░░░░░░░░░░░░░░░░░░░█
████████▄█▄▀▄░░▀████▄░░░░░░░░░░░░░░░░░░░░░░░░█
███░░░▄▀▀▀░░░░▄░▄████▀░░░░░▄▄▄▄▄▄▄░░░░░░░░░░░█
██░░░░░▀▄░░░░░░▀████▄░░░░▄█▄███████▄▄▄░▄▄░░░░█
█░░░░░░░░▀░░░░▄████▀░░░░░░░█████▀███████▀▄░░░█
█░░░░░░░░░░░▄██████░░░░░░░░░███▄▄▄░░▀░▀████▄░█
█░░░░░░░░░░░█████▀░░░░░░░░░░▀░▀█▀░▀▄▄▄░░░▄██░█
█░░░░░░░░░░░▄███▄░░░░░░░░░░░▀▄░▄░░░░░▀▀▄░░▀░░█
█░░░░░░░░░▄░▄██▀░░░░░░░░░░░░░░░░░░░░░░░░▀░░░░█
██▄░░░░░▄▄██▄█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█
████▄░▄▄█████▀░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▄█
████▄░▄███▄███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█
████▄▀░▀░▀▄▀██▄░░░░░▄░░░░░░░░░░░░░░░░░░░░░░▄██
████▄█▄░▄░▄█▄░░▀░░░▀▄█▄░░░░░░░░░░░░░░░░░░▄████
█████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░▄█████
████████████▄░░█▄░░░▄░░░░░░░░░░░░░░░░░░▄██████
██████████▄▄▄░░▀▀▀██▄░░░░░░░░░░░░░░░░░▄███████
████▄░░░░▀████▄▀▄▄▄▀▄▀▄▄▄░░░░░░░░░░░░█████████
██████▄░░░░░░▀░▀░▀███▄███▄░░░░░░░░░▄██████████
████████▄▄▄▄▄░░▀░▀░░░░░░░▀████▄░░░▄███████████
█████████▀░░░░░░░░░░░░░░░░░███░░▄█████████████
█████▀░░░░░░░░░░░░░░░░░░░░░▀█▀░▄██████████████
███▀░░░░░░░░░░▄██████▄░░░░░░░▀▄███████████████
███▄░░░░▄█████████████▄▄▄▄▄░▄▄████████████████
█████████████████▀▀▀░░░░░░░░░░░░▀█████████████
███████████▀░░░░░░░░░░░░░░░░░░░░░░▀███████████
█████████▀░░░░░░░░░░░░░░░░░▄▄░▄▄▄░░░██████████
███████▀░░░░░░░░░░▄▄▄▄▄▄▄▄▀░░░░░░░░░░█████████
█████▀░░░░░░░░▄▄▄██▀░░░░░░░░░░░░░░░░░░░███████
███▀░░░░░░░▄▄██▀░░░░░░░░░░░░░▄▄▄▄▄▄░░░░░██████
███░░░░░░▄▄█░░░░░░░░░░░░▄▄▄▄▄▀░░░░░░░░░░░█████
██░░░░░░▄██▀░░░░▄▄▄▄▄▀▀▀░░░░░░░░░░░░░░░░░░████
█░░░░░▄███░░░░▄████░░░░░░░░░░░░░░░░░░░░░░░░███
░░░░░▄███▀░░░▄███▀░░░░░░░░░░▄░░░░░░░░░░░░░░░██
░░░░▄███░░░░███▀░░░░░░▄██████▀░░░░░░░░░░░░░░░█

Don't worry about a thing,
'Cause every little thing gonna be all right.

Bob Marley
</pre>`),
	}

	if c.ViewUser.Language == "ru" {
		c.StaticPage.Title = "Помощь"
	}

	if c.StaticPage.Title == "" {
		http.NotFound(w, r.Request)
		return
	}

	util.RenderTemplate(w, "help", c)
}

func (c *Context) HelpItem(w web.ResponseWriter, r *web.Request) {
	lang := c.ViewUser.Language
	if lang == "" {
		lang = "en"
	}
	pages, err := loadStaticPages(fmt.Sprintf("./static/%s/", lang))
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.StaticPages = pages

	for _, page := range pages {
		if page.Filename == r.PathParams["filename"] {
			c.StaticPage = page
		}
	}

	util.RenderTemplate(w, "help", c)
}
