package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) SellerMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	user, _ := FindUserByUsername(r.PathParams["store"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}

	seller := Seller{user}
	c.ViewSeller = seller.ViewSeller(c.ViewUser.Language)

	c.CanEdit = (c.ViewUser.Uuid == c.ViewSeller.Uuid) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff
	if !c.CanEdit {
		http.NotFound(w, r.Request)
		return
	}

	next(w, r)
}

func (c *Context) SellerItemMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	c.CanEdit = c.ViewUser.IsAdmin || c.ViewUser.IsTrustedSeller

	if r.PathParams["item"] != "" && r.PathParams["item"] != "new" {
		item, _ := FindItemByUuid(r.PathParams["item"])
		if item == nil || (item.User.Uuid != c.ViewUser.Uuid && !c.CanEdit) {
			http.NotFound(w, r.Request)
			return
		}
		c.Item = *item
		c.ViewItem = c.Item.ViewItem(c.ViewUser.Language)
	}

	if r.PathParams["item"] == "new" {
		items := FindItemsForSeller(c.ViewSeller.Uuid)
		numberOfItems := len(items)
		if numberOfItems >= 5 && !c.ViewSeller.Premium {
			http.Redirect(w, r.Request, "/free_restrictions", 302)
			return
		}
	}

	store, _ := FindUserByUsername(r.PathParams["store"])
	if store.Username != c.ViewUser.Username && !c.CanEdit {
		http.NotFound(w, r.Request)
		return
	}

	next(w, r)
}

func (c *Context) SellerItemPackageMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if r.PathParams["package"] != "" && r.PathParams["package"] != "new" {
		itemPackage, _ := FindPackageByUuid(r.PathParams["package"])
		if itemPackage != nil {
			if itemPackage.ItemUuid != c.Item.Uuid {
				http.NotFound(w, r.Request)
				return
			}
			c.Package = *itemPackage
			c.ViewPackage = itemPackage.ViewPackage()
		}
	}
	next(w, r)
}

func (c *Context) VendorMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	seller, _ := FindSellerByUsername(r.PathParams["store"])
	if seller == nil || seller.Banned {
		http.NotFound(w, r.Request)
		return
	}
	reviews, _ := FindRatingReviewsBySellerUuid(seller.Uuid)
	seller.RatingReviews = reviews
	c.ViewSeller = seller.ViewSeller(c.ViewUser.Language)
	c.ViewItems = Items(seller.Items).ViewItems(c.ViewUser.Language)
	c.CanEdit = (c.ViewUser.Uuid == c.ViewSeller.Uuid) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff

	thread, err := GetStoreThread(*seller)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)

	next(w, r)
}

func (c *Context) VendorItemMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	item, _ := FindItemByUuid(r.PathParams["item"])
	if item == nil || item.UserUuid != c.ViewSeller.Uuid {
		http.NotFound(w, r.Request)
		return
	}
	c.Item = *item
	c.ViewItem = item.ViewItem(c.ViewUser.Language)
	next(w, r)
}
