package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) MessagesMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if recieverUsername := r.PathParams["username"]; recieverUsername != "" {
		reciever, _ := FindUserByUsername(recieverUsername)
		if reciever == nil {
			http.NotFound(w, r.Request)
			return
		}
		thread, err := GetPrivateThread(*c.ViewUser.User, *reciever, "", r.Method == "POST")
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.Thread = *thread
		c.ViewThread = thread.ViewThread(c.ViewUser.User.Language, c.ViewUser.User)
	}

	c.ViewThreads = FindPrivateThreads(*c.ViewUser.User).ViewThreads(c.ViewUser.User.Language, c.ViewUser.User)

	next(w, r)
}

func (c *Context) MessageStatsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.Uuid == "" {
		next(w, r)
		return
	}

	c.NumberOfPrivateMessages = CountPrivateMessages(*c.ViewUser.User)
	c.NumberOfUnreadPrivateMessages = CountUndreadPrivateMessages(*c.ViewUser.User)
	// c.NumberOfUnreadSupportMessages = CountUndreadSupportMessages(*c.ViewUser.User)
	c.NumberOfSupportMessages = CountSupportTicketsForUser(*c.ViewUser.User)

	next(w, r)
}
