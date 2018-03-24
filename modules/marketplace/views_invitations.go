package marketplace

import (
	"github.com/gocraft/web"
	"net/http"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ListInvitations(w web.ResponseWriter, r *web.Request) {
	invitatations := FindInvitationsByInviterUuid(c.ViewUser.Uuid)
	seller := Seller{c.ViewUser.User}
	c.ViewSeller = seller.ViewSeller(c.ViewUser.Language)
	c.Invitations = invitatations
	util.RenderTemplate(w, "invitations/list", c)
}

func (c *Context) ShowInvitation(w web.ResponseWriter, r *web.Request) {
	uuid := r.PathParams["uuid"]
	invite, err := FindInvitationByUuid(uuid)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.Invitation = *invite
	util.RenderTemplate(w, "invitations/show", c)
}

func (c *Context) DeleteInvitation(w web.ResponseWriter, r *web.Request) {
	uuid := r.PathParams["uuid"]
	if uuid != "new" {
		invite, err := FindInvitationByUuid(uuid)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		invite.Remove()
	}

	http.Redirect(w, r.Request, "/invitations", 302)
}

func (c *Context) EditInvitation(w web.ResponseWriter, r *web.Request) {

	seller := Seller{c.ViewUser.User}
	c.ViewSeller = seller.ViewSeller(c.ViewUser.Language)
	uuid := r.PathParams["uuid"]
	if uuid != "new" {
		invite, err := FindInvitationByUuid(uuid)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.Invitation = *invite
	}
	util.RenderTemplate(w, "invitations/edit", c)
}

func (c *Context) SaveInvitation(w web.ResponseWriter, r *web.Request) {

	if r.PathParams["uuid"] == "new" {
		c.Invitation.Uuid = util.GenerateUuid()
	} else if r.PathParams["uuid"] != "" {
		c.Invitation.Uuid = r.PathParams["uuid"]
	}

	c.Invitation.Username = r.FormValue("username")
	c.Invitation.InvitationText = r.FormValue("invitation_text")
	c.Invitation.InviterUuid = c.ViewUser.Uuid

	validationError := c.Invitation.Save()
	if validationError != nil {
		c.Error = validationError.Error()
		c.EditInvitation(w, r)
		return
	}

	http.Redirect(w, r.Request, "/invitations/"+c.Invitation.Uuid, 302)
}
