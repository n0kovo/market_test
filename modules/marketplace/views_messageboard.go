package marketplace

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"github.com/n0kovo/market_test/modules/util"
)

func (c *Context) ListGeneralThreads(w web.ResponseWriter, r *web.Request) {

	if len(r.URL.Query()["section"]) > 0 {
		section := r.URL.Query()["section"][0]
		sectionID, _ := strconv.ParseInt(section, 10, 64)
		c.SelectedSectionID = int(sectionID)
	} else {
		c.SelectedSectionID = 1
	}

	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page) - 1
	}

	numberOfThreadsPerPage := 50

	messagebordThreads := FindMessageboardThreadsForUserUuid(c.SelectedSectionID, c.Page, numberOfThreadsPerPage, c.ViewUser.Uuid)
	c.ViewMessageboardThreads = MessageboardThreads(messagebordThreads).ViewMessageboardThreads(c.ViewUser.Language)
	c.MessageboardSections = FindParentMessageboardSections()

	numberOfThreads := float64(CountMessageboardThreads(c.SelectedSectionID))
	c.NumberOfPages = int(math.Ceil(numberOfThreads / float64(numberOfThreadsPerPage)))

	// paging
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}
	c.Page += 1

	util.RenderTemplate(w, "board/threads", c)
}

func (c *Context) ListSellerThreads(w web.ResponseWriter, r *web.Request) {

	c.ViewThreads = FindSellerThreads().ViewThreads(c.ViewUser.Language, c.ViewUser.User)
	c.MessageboardSections = FindAllMessageboardSections()

	util.RenderTemplate(w, "board/seller_threads", c)
}

func (c *Context) ShowThread(w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.Uuid == "" {
		redirectUrl := "/auth/register"
		http.Redirect(w, r.Request, redirectUrl, 302)
		return
	}

	thread, err := GetMessageboardThread(r.PathParams["uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	c.NumberOfPages = int(math.Ceil(float64(len(c.ViewThread.Messages)) / 50.0))

	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page) - 1
	}
	// paging
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}

	c.ViewThread.Messages = c.ViewThread.Messages[c.Page*50 : int(math.Min(float64(len(c.ViewThread.Messages)), float64(c.Page*50+50)))]
	c.Page = c.Page + 1
	c.CaptchaId = captcha.New()

	// c.ViewThreads = FindMessageboardThreads(c.SelectedSectionID).ViewThreads(c.ViewUser.Language, c.ViewUser.User)
	c.SelectedSection = c.ViewThread.Section

	UpdateThreadPerusalStatus(thread.Uuid, c.ViewUser.Uuid)
	c.MessageboardSections = FindParentMessageboardSections()

	util.RenderTemplate(w, "board/thread", c)
}

func (c *Context) MessageImage(w web.ResponseWriter, r *web.Request) {
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	util.ServeImage(r.PathParams["uuid"], size, w, r)
}

func (c *Context) DeleteThread(w web.ResponseWriter, r *web.Request) {
	thread, err := FindThreadByUuid(r.PathParams["uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	if thread.SenderUserUuid == c.ViewUser.Uuid || c.ViewUser.IsAdmin || c.ViewUser.IsStaff {
		thread.Remove()
	}
	http.Redirect(w, r.Request, "/board/", 302)
}

func (c *Context) EditThread(w web.ResponseWriter, r *web.Request) {
	var editThread bool
	if r.PathParams["uuid"] != "" {
		thread, err := GetMessageboardThread(r.PathParams["uuid"])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		editThread = true
	}
	c.CaptchaId = captcha.New()
	c.MessageboardSections = FindAllMessageboardSections()

	if editThread {
		util.RenderTemplate(w, "board/thread_edit", c)
	} else {
		util.RenderTemplate(w, "board/thread_new", c)
	}
}

func (c *Context) EditThreadPOST(w web.ResponseWriter, r *web.Request) {
	// vars
	var (
		thread      *Thread
		isNewThread bool
		err         error
	)

	// captcha
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.EditThread(w, r)
		return
	}

	// new or existing thread
	if r.PathParams["uuid"] != "" {
		thread, err = GetMessageboardThread(r.PathParams["uuid"])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
	} else {
		thread, err = CreateThread(
			"messageboard",
			"",
			r.FormValue("title"),
			r.FormValue("text"),
			c.ViewUser.User,
			nil,
			true,
		)
		if err != nil {
			c.Error = err.Error()
			c.EditThread(w, r)
			return
		}
		isNewThread = true
	}

	// section
	secId, err := strconv.ParseInt(r.FormValue("section_id"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.EditThread(w, r)
		return
	}
	section, err := FindMessageboardSectionByID(int(secId))
	if err != nil {
		c.Error = err.Error()
		c.EditThread(w, r)
		return
	}

	// set title, text and section
	thread.Title = r.FormValue("title")
	thread.Text = r.FormValue("text")
	thread.MessageboardSectionID = section.ID
	thread.Save()

	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	err = thread.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.EditThread(w, r)
		return
	}

	// feed actions
	if isNewThread {
		CreateFeedItem(c.ViewUser.Uuid, "new_thread", "created new thread", thread.Uuid)
	}

	// redirect
	http.Redirect(w, r.Request, fmt.Sprintf("/board/?section=%d", thread.MessageboardSectionID), 302)
}

func (c *Context) ReplyToThread(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ShowThread(w, r)
		return
	}
	thread, err := GetMessageboardThread(r.FormValue("thread_uuid"))
	if err != nil {
		c.Error = err.Error()
		c.ShowThread(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowThread(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ShowThread(w, r)
		return
	}

	CreateFeedItem(c.ViewUser.Uuid, "new_thread_reply", "replied in thread", message.Uuid)
	c.ShowThread(w, r)
}

func (c *Context) ListPrivateMessages(w web.ResponseWriter, r *web.Request) {
	c.ViewSeller = Seller{c.ViewUser.User}.ViewSeller(c.ViewUser.Language) //@

	util.RenderTemplate(w, "board/messages", c)
}

func (c *Context) ShowPrivateMessage(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	c.ViewSeller = Seller{c.ViewUser.User}.ViewSeller(c.ViewUser.Language)
	c.SelectedSection = c.ViewThread.Uuid

	for _, m := range c.ViewThread.Messages {
		if m.RecieverUserUuid == c.ViewUser.Uuid && !m.IsReadByReciever {
			m.IsReadByReciever = !m.IsReadByReciever
			m.Save()
		}
	}
	UpdateThreadPerusalStatus(c.ViewThread.Uuid, c.ViewUser.Uuid)

	// hack to make displayed thread viewd as read
	for i, _ := range c.ViewThreads {
		th := c.ViewThreads[i]
		if th.Uuid == c.ViewThread.Uuid {
			c.ViewThreads[i].IsRead = true
			break
		}
	}

	util.RenderTemplate(w, "board/message", c)
}

func (c *Context) ShowPrivateMessagePOST(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := captcha.VerifyString(
		r.FormValue("captcha_id"),
		r.FormValue("captcha"),
	)
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ShowPrivateMessage(w, r)
		return
	}

	message, err := CreateMessage(r.FormValue("text"), c.Thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ShowPrivateMessage(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ShowPrivateMessage(w, r)
		return
	}

	c.MessagesMiddleware(w, r, c.ShowPrivateMessage)
}
