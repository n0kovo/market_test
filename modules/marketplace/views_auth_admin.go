package marketplace

import (
	"github.com/gocraft/web"
	"math"
	"net/http"
	"github.com/n0kovo/market_test/modules/util"
	"strconv"
	"time"
)

func (c *Context) AdminUsers(w web.ResponseWriter, r *web.Request) {
	numberOfUsers := CountUsers(nil)

	// paging
	pageSize := 50
	selectedPage := 0
	if len(r.URL.Query()["page"]) > 0 {
		selectedPageStr := r.URL.Query()["page"][0]
		page, err := strconv.Atoi(selectedPageStr)
		if err == nil {
			selectedPage = page - 1
		}
	}
	numberOfPages := int(math.Ceil(float64(numberOfUsers) / float64(pageSize)))
	for i := 0; i < numberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}
	// Sort By
	if len(r.URL.Query()["sortby"]) > 0 {
		c.SortBy = r.URL.Query()["sortby"][0]
	}
	// Query
	if len(r.URL.Query()["query"]) > 0 {
		c.Query = r.URL.Query()["query"][0]
	}
	usersPage := GetExtendedUsersPage(selectedPage, pageSize, c.SortBy, c.Query)
	c.ExtendedUsers = usersPage

	// Stats
	c.NumberOfPages = numberOfPages
	c.SelectedPage = selectedPage + 1
	c.NumberOfUsers = numberOfUsers
	c.NumberOfNewUsers = CountUsersRegistredAfter(time.Now().AddDate(0, 0, -1))
	c.NumberOfActiveUsers = CountUsersActiveAfter(time.Now().AddDate(0, 0, -1))
	c.NumberOfWeeklyActiveUsers = CountUsersActiveAfter(time.Now().AddDate(0, 0, -7))
	c.NumberOfMonthlyActiveUsers = CountUsersActiveAfter(time.Now().AddDate(0, -1, 0))
	d, _ := time.ParseDuration("-10m")
	c.NumberOfOnlineUsers = CountUsersActiveAfter(time.Now().Add(d))

	util.RenderTemplate(w, "auth/admin/users", c)
}

func (c *Context) LoginAsUser(w web.ResponseWriter, r *web.Request) {
	if !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}
	username := r.PathParams["user"]
	user, _ := FindUserByUsername(username)
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	c.Login(*user, w, r)
}

func (c *Context) GrantPremium(w web.ResponseWriter, r *web.Request) {
	username := r.PathParams["user"]
	user, _ := FindUserByUsername(username)
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.Premium = !user.Premium
	user.Save()

	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) GrantPremiumPlus(w web.ResponseWriter, r *web.Request) {
	username := r.PathParams["user"]
	user, _ := FindUserByUsername(username)
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.PremiumPlus = !user.PremiumPlus
	user.Save()

	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) GrantSeller(w web.ResponseWriter, r *web.Request) {
	user, _ := FindUserByUsername(r.PathParams["user"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.IsSeller = !user.IsSeller
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) GrantTrustedSeller(w web.ResponseWriter, r *web.Request) {
	user, _ := FindUserByUsername(r.PathParams["user"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.IsTrustedSeller = !user.IsTrustedSeller
	user.TrusteeUuid = c.ViewUser.Uuid
	user.HasRequestedVerification = false
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) GrantStaff(w web.ResponseWriter, r *web.Request) {
	if !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}
	user, _ := FindUserByUsername(r.PathParams["user"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.IsStaff = !user.IsStaff
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) GrantTester(w web.ResponseWriter, r *web.Request) {
	user, _ := FindUserByUsername(r.PathParams["user"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.IsTester = !user.IsTester
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) BanUser(w web.ResponseWriter, r *web.Request) {
	username := r.PathParams["user"]
	user, _ := FindUserByUsername(username)
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.Banned = !user.Banned
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}

func (c *Context) MarkPossibleScammer(w web.ResponseWriter, r *web.Request) {
	username := r.PathParams["user"]
	user, _ := FindUserByUsername(username)
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	user.PossibleScammer = !user.PossibleScammer
	user.Save()
	http.Redirect(w, r.Request, r.Referer(), 302)
}
