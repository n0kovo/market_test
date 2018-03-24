package marketplace

import (
	"math"
	"net/http"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewStaffListUsers(w web.ResponseWriter, r *web.Request) {

	var (
		users         []ExtendedUser
		err           error
		numberOfUsers int
		pageSize      int = 50
		selectedPage  int = 0
	)

	if len(r.URL.Query()["page"]) > 0 {
		selectedPageStr := r.URL.Query()["page"][0]
		page, err := strconv.Atoi(selectedPageStr)
		if err == nil {
			selectedPage = page - 1
		}
	}

	c.SelectedSection = "uncontacted"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	if c.SelectedSection == "uncontacted" {
		users, err = FindUncontactedUsers(selectedPage, pageSize)
		numberOfUsers = CountUncontactedUsers()

	} else if c.SelectedSection == "contacted" {
		users, err = FindUsersContactedByStaff(selectedPage, pageSize)
		numberOfUsers = CountUsersContactedByStaff(c.ViewUser.Uuid)
	}

	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	numberOfPages := int(math.Ceil(float64(numberOfUsers) / float64(pageSize)))
	for i := 0; i < numberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}

	c.ViewExtendedUsers = ExtendedUsers(users).ViewExtendedUsers(c.Language)
	c.SelectedPage = selectedPage + 1

	util.RenderTemplate(w, "staff/users_support_general", c)
}

func (c *Context) ViewStaffListSupportTickets(w web.ResponseWriter, r *web.Request) {
	var (
		err error
	)

	c.SelectedSection = "all"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	tickets, err := FindSupportTicketsByStatus(c.SelectedSection)
	if err != nil {
		panic(err)
		http.NotFound(w, r.Request)
		return
	}

	c.ViewSupportTickets = tickets.ViewSupportTickets(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/users_support_tickets", c)
}

func (c *Context) ViewStaffUserFinance(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.ViewSeller(c.ViewUser.Language)

	c.ViewSeller.BitcoinWallets = c.ViewSeller.FindUserBitcoinWallets()
	c.ViewSeller.EthereumWallets = c.ViewSeller.FindUserEthereumWallets()
	c.ViewSeller.BitcoinCashWallets = c.ViewSeller.FindUserBitcoinCashWallets()

	for _, w := range c.ViewSeller.BitcoinWallets {
		w.UpdateBalance(false)
	}

	for _, w := range c.ViewSeller.EthereumWallets {
		w.UpdateBalance(false)
	}

	for _, w := range c.ViewSeller.BitcoinCashWallets {
		w.UpdateBalance(false)
	}

	c.ViewSeller.BitcoinBalance = c.ViewSeller.BitcoinWallets.Balance()
	c.ViewSeller.EthereumBalance.Balance = c.ViewSeller.EthereumWallets.Balance().Balance
	c.ViewSeller.BitcoinCashBalance = c.ViewSeller.BitcoinCashWallets.Balance()

	c.ViewSeller.BitcoinWallet = c.ViewSeller.BitcoinWallets[0]
	c.ViewSeller.EthereumWallet = c.ViewSeller.EthereumWallets[0]
	c.ViewSeller.BitcoinCashWallet = c.ViewSeller.BitcoinCashWallets[0]

	c.UserSettingsHistory = SettingsChangeHistoryByUser(user.Uuid)

	util.RenderTemplate(w, "staff/users_user_finance", c)
}

func (c *Context) ViewStaffUserTickets(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.ViewSeller(c.ViewUser.Language)

	tickets, err := FindSupportTicketsForUser(*user)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSupportTickets = tickets.ViewSupportTickets(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/users_user_tickets", c)
}

func (c *Context) ViewStaffUserPayments(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.
		ViewSeller(c.ViewUser.Language)

	c.ViewCurrentTransactionStatuses = FindCurrentTransactionStatuses(
		user.Uuid, c.SelectedStatus, false, 0, 100).
		ViewCurrentTransactionStatuses(c.ViewUser.Language)

	util.RenderTemplate(w, "staff/users_user_payments", c)
}

func (c *Context) ViewStaffUserAdminActions(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.ViewSeller(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/users_user_admin_actions", c)
}

// Support chat

func (c *Context) ViewStaffGeneralSupportThreadGET(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.ViewSeller(c.ViewUser.Language)

	c.CaptchaId = captcha.New()
	thread, err := GetSupportThread(*user, false)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)

	util.RenderTemplate(w, "staff/users_user_support", c)
}

func (c *Context) ViewStaffGeneralSupportThreadPOST(w web.ResponseWriter, r *web.Request) {

	user, _ := FindUserByUsername(r.PathParams["username"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSeller = Seller{user}.ViewSeller(c.ViewUser.Language)

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewStaffGeneralSupportThreadGET(w, r)
		return
	}

	thread, err := GetSupportThread(*user, true)
	if err != nil {
		c.Error = err.Error()
		c.ViewStaffGeneralSupportThreadGET(w, r)
		return
	}

	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewStaffGeneralSupportThreadGET(w, r)
		return
	}

	err = message.AddImage(r)

	_, _, err = r.FormFile("image")
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewStaffGeneralSupportThreadGET(w, r)
		return
	}

	if user.SupporterUuid == "" && thread.SenderUserUuid != c.ViewUser.Uuid {
		user.SupporterUuid = c.ViewUser.Uuid
		user.Save()
	}

	CreateFeedItem(c.ViewUser.Uuid, "staff_support_reply", "replied in support thread", user.Uuid)
	EventNewSupportMessage(message)

	c.ViewStaffGeneralSupportThreadGET(w, r)
}
