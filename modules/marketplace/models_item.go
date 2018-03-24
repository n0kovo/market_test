package marketplace

import (
	"encoding/base32"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/russross/blackfriday"
)

/*
	Models
*/

type Item struct {
	Uuid           string `json:"uuid" gorm:"primary_key"`
	Name           string `json:"name"`
	Description    string `json:"description" sql:"size:4096"`
	ItemCategoryID int    `json:"category_id" gorm:"index"`
	UserUuid       string `json:"user_uuid" gorm:"index"`
	IsPromoted     bool   `json:"is_promoted"  gorm:"index"`
	NumberOfSales  int    `json:"number_of_sales" gorm:"index"`
	NumberOfViews  int    `json:"number_of_views" gorm:"index"`

	User          User           `json:"-"`
	ReviewerUser  User           `json:"-" gorm:"ForeignKey:ReviewedByUserUuid"`
	ItemCategory  ItemCategory   `json:"-"`
	Packages      []Package      `json:"-"`
	RatingReviews []RatingReview `json:"-"`

	ReviewedByUserUuid string
	ReviewedAt         *time.Time

	CreatedAt *time.Time `json:"created_at" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type Items []Item

type ViewItems []ViewItem
type ViewItemsByPrice []ViewItem

func (a ViewItemsByPrice) Len() int      { return len(a) }
func (a ViewItemsByPrice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ViewItemsByPrice) Less(i, j int) bool {
	if a[i].User.PremiumPlus != a[j].User.PremiumPlus { // by premium status
		return a[i].User.PremiumPlus
	} else if a[i].User.Premium != a[j].User.Premium { // by premium status
		return a[i].User.Premium
	} else if a[i].ScoreFloat != a[j].ScoreFloat { // by score
		return a[i].ScoreFloat > a[j].ScoreFloat
	} else { // by price
		return a[i].MedPriceBTCFloat < a[j].MedPriceBTCFloat
	}
}

type ViewItemsByPopularity []ViewItem

func (s ViewItemsByPopularity) Len() int {
	return len(s)
}

func (s ViewItemsByPopularity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ViewItemsByPopularity) Less(i, j int) bool {
	if s[i].User.PremiumPlus != s[j].User.PremiumPlus { // by premium status
		return s[i].User.PremiumPlus
	} else if s[i].User.Premium != s[j].User.Premium { // by premium status
		return s[i].User.Premium
	} else {
		return s[i].NumberOfSales > s[j].NumberOfSales
	}
}

type ViewItemsByDateAdded []ViewItem

func (s ViewItemsByDateAdded) Len() int {
	return len(s)
}

func (s ViewItemsByDateAdded) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ViewItemsByDateAdded) Less(i, j int) bool {
	if s[i].CreatedAt != nil && s[j].CreatedAt != nil {
		return s[i].CreatedAt.After(*s[j].CreatedAt)
	}
	return true
}

type ViewItemsByDateLoggedIn []ViewItem

func (s ViewItemsByDateLoggedIn) Len() int {
	return len(s)
}

func (s ViewItemsByDateLoggedIn) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ViewItemsByDateLoggedIn) Less(i, j int) bool {
	if s[i].User.LastLoginDate == nil || s[j].User.LastLoginDate == nil {
		return false
	}
	return s[i].User.LastLoginDate.After(*(s[j].User.LastLoginDate))
}

type ViewItemsByRating []ViewItem

func (s ViewItemsByRating) Len() int {
	return len(s)
}

func (s ViewItemsByRating) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ViewItemsByRating) Less(i, j int) bool {
	if s[i].User.PremiumPlus != s[j].User.PremiumPlus { // by premium status
		return s[i].User.PremiumPlus
	} else if s[i].User.Premium != s[j].User.Premium { // by premium status
		return s[i].User.Premium
	} else {
		return s[i].ScoreFloat > s[j].ScoreFloat
	}
}

/*
	Model Methods
*/

func (itms Items) Where(fnc func(Item) bool) Items {
	filtered := Items{}
	for _, i := range itms {
		if fnc(i) {
			filtered = append(filtered, i)
		}
	}
	return filtered
}

func (item Item) ImageBase32() string {
	content, err := ioutil.ReadFile("./data/images/" + item.Uuid + "_200x200.jpeg")
	if err != nil {
		return ""
	}
	return "image/jpeg;base32," + base32.StdEncoding.EncodeToString(content)
}

func (i Item) Score() float64 {
	score := float64(0.0)
	for _, r := range i.RatingReviews {
		score += float64(r.ItemScore)
	}
	if len(i.RatingReviews) > 0 {
		score /= float64(len(i.RatingReviews))
	}
	return math.Ceil(score*100) / float64(100.0)
}

func (i *Item) Validate() error {
	if i.Name == "" {
		return errors.New("Name is not valid")
	}
	if i.Description == "" {
		return errors.New("Description is not valid")
	}
	if i.UserUuid == "" {
		return errors.New("UserUuid is not valid")
	}
	if i.User.Username == "" {
		if i.UserUuid != "" {
			user, err := FindUserByUuid(i.UserUuid, false)
			if err != nil {
				return err
			}
			i.User = *user
		} else {
			return errors.New("No such seller")
		}
	}

	return nil
}

/*
	Database methods
*/

func (i Item) Remove() error {
	go BleveIndex.Delete(i.Uuid)
	return database.Delete(&i).Error
}

func (itm Item) Save() error {
	err := itm.Validate()
	if err != nil {
		return err
	}
	go itm.Index()
	return itm.SaveToDatabase()
}

func (itm Item) SaveToDatabase() error {
	if existing, _ := FindItemByUuid(itm.Uuid); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

func (i Item) Index() error {
	return BleveIndex.Index(i.Uuid, i)
}

/*
	Relations
*/

func (item Item) PackagesWithoutReservation() Packages {
	predicate := func(a Package) bool {
		return a.Reservation() == nil
	}
	for i, _ := range item.Packages {
		item.Packages[i].Item = item
	}
	return Packages(item.Packages).Where(predicate)
}

/*
	Queries
*/

func CountItems() int {
	var count int
	database.Table("items").Count(&count)
	return count
}

func GetAllItems() Items {
	var items []Item
	database.Unscoped().Find(&items)
	return Items(items)
}

func FindItemByUuid(uuid string) (*Item, error) {
	var item Item
	err := database.
		Preload("Packages").
		Preload("Packages.PackagePrice").
		Preload("Packages.GeoCity").
		Preload("Packages.GeoCountryFrom").
		Preload("Packages.GeoCountryTo").
		Preload("User").
		Preload("ReviewerUser").
		Preload("RatingReviews").
		Preload("RatingReviews.User").
		Preload("ItemCategory").
		First(&item, "uuid = ?", uuid).
		Error
	if err != nil {
		return nil, err
	}

	reviews, _ := FindRatingReviewsBySellerUuid(item.User.Uuid)
	item.User.RatingReviews = reviews

	return &item, err
}

// FindActiveItems returns items with reservations avaiable for reservation.
func FindActiveItems() Items {
	var items []Item
	database.
		Joins("join users on users.uuid=items.user_uuid").
		Where("users.banned=0 and users.is_seller=true and users.vacation_mode=false").
		Find(&items)
	return Items(items)
}

// FindActiveItems returns items with reservations avaiable for reservation.
func FindActiveItemsForSellerUuid(sellerUuid string) Items {
	var items []Item
	database.
		Joins("join users on users.uuid=items.user_uuid").
		Where("users.banned=0 and users.is_seller=true and users.vacation_mode=false and users.uuid=?", sellerUuid).
		Find(&items)
	return Items(items)
}

func FindItemsForSeller(uuid string) Items {
	var items []Item
	database.
		Where(&Item{UserUuid: uuid}).
		Preload("User").
		Preload("RatingReviews").
		Preload("ItemCategory").
		Preload("Packages").
		Preload("Packages.PackagePrice").
		Preload("Packages.GeoCity").
		Preload("Packages.GeoCountryFrom").
		Preload("Packages.GeoCountryTo").
		Find(&items)
	return Items(items)
}

func CountUnreviewedItems() int {
	var count int
	database.
		Table("items").
		Joins("join users on users.uuid=items.user_uuid").
		Where("users.banned=false").
		Where("reviewed_by_user_uuid is null or reviewed_by_user_uuid=''").
		Count(&count)
	return count
}

func FindUnreviewedItems(page, pageSize int) Items {
	var items []Item
	database.
		Table("items").
		Joins("join users on users.uuid=items.user_uuid").
		Where("users.banned=false").
		Where("reviewed_by_user_uuid is null or reviewed_by_user_uuid=''").
		Limit(pageSize).
		Offset(page).
		Preload("User").
		Preload("ItemCategory").
		Preload("Packages").
		Preload("Packages.PackagePrice").
		Preload("Packages.GeoCity").
		Preload("Packages.GeoCountryFrom").
		Preload("Packages.GeoCountryTo").
		Order("updated_at DESC").
		Find(&items)
	return Items(items)
}

func FindTopSellerItems() Items {
	var items []Item
	database.
		Joins("join users on users.uuid=items.user_uuid").
		Where(`
			users.banned=false and 
			users.is_seller=true and 
			users.vacation_mode=false and 
			users.last_login_date > ?`,
			time.Now().AddDate(0, 0, -MARKETPLACE_SETTINGS.CooloffPeriod),
		).
		Preload("User").
		Preload("RatingReviews").
		Preload("ItemCategory").
		Preload("Packages").
		Preload("Packages.PackagePrice").
		Preload("Packages.GeoCity").
		Preload("Packages.GeoCountryFrom").
		Preload("Packages.GeoCountryTo").
		Order("number_of_sales desc").
		Limit(20).
		Find(&items)
	return Items(items)
}

/*
	View models
*/

type ViewItem struct {
	*Item
	DescriptionHTML      template.HTML
	ShortDescriptionHTML template.HTML
	GroupPackages        []GroupPackage

	MaxPriceBTCStr   string
	MaxPriceBTCFloat float64
	MinPriceBTCStr   string
	MinPriceBTCFloat float64

	MaxPriceBCHStr   string
	MaxPriceBCHFloat float64
	MinPriceBCHStr   string
	MinPriceBCHFloat float64

	MaxPriceETHStr   string
	MaxPriceETHFloat float64
	MinPriceETHStr   string
	MinPriceETHFloat float64

	MaxPriceUSDStr   string
	MaxPriceUSDFloat float64
	MinPriceUSDStr   string
	MinPriceUSDFloat float64

	MaxPriceEURStr   string
	MaxPriceEURFloat float64
	MinPriceEURStr   string
	MinPriceEURFloat float64

	MaxPriceGBPStr   string
	MaxPriceGBPFloat float64
	MinPriceGBPStr   string
	MinPriceGBPFloat float64

	MaxPriceAUDStr   string
	MaxPriceAUDFloat float64
	MinPriceAUDStr   string
	MinPriceAUDFloat float64

	MaxPriceRUBStr   string
	MaxPriceRUBFloat float64
	MinPriceRUBStr   string
	MinPriceRUBFloat float64

	MedPriceBTCFloat float64

	ViewSeller        ViewSeller
	ViewPackages      []ViewPackage
	ViewRatingReviews []ViewRatingReview

	ScoreStr   string
	ScoreFloat float64

	Premium bool
}

func (item Item) ViewItem(lang string) ViewItem {
	itemPackages := Packages(item.Packages)
	score := item.Score()

	pckgs, _ := Packages(item.Packages).GroupsTable().GetPackagesPage(0, 5)

	vi := ViewItem{
		Item:            &item,
		DescriptionHTML: template.HTML(userHtmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(item.Description))))),
		ShortDescriptionHTML: template.HTML(
			userHtmlPolicy.Sanitize(
				string(
					blackfriday.MarkdownCommon(
						[]byte(
							item.Description[0:int(
								math.Min(
									float64(len(item.Description)),
									float64(1024),
								),
							)],
						),
					),
				),
			),
		),

		ViewSeller:   (Seller{&item.User}).ViewSeller(lang),
		ViewPackages: itemPackages.ViewPackages(),

		MinPriceBTCFloat: itemPackages.MinPrice("BTC"),
		MaxPriceBTCFloat: itemPackages.MaxPrice("BTC"),
		MinPriceBTCStr:   humanize.Ftoa(itemPackages.MinPrice("BTC")),
		MaxPriceBTCStr:   humanize.Ftoa(itemPackages.MaxPrice("BTC")),

		MinPriceBCHFloat: itemPackages.MinPrice("BCH"),
		MaxPriceBCHFloat: itemPackages.MaxPrice("BCH"),
		MinPriceBCHStr:   humanize.Ftoa(itemPackages.MinPrice("BCH")),
		MaxPriceBCHStr:   humanize.Ftoa(itemPackages.MaxPrice("BCH")),

		MinPriceETHFloat: itemPackages.MinPrice("ETH"),
		MaxPriceETHFloat: itemPackages.MaxPrice("ETH"),
		MinPriceETHStr:   humanize.Ftoa(itemPackages.MinPrice("ETH")),
		MaxPriceETHStr:   humanize.Ftoa(itemPackages.MaxPrice("ETH")),

		MinPriceUSDFloat: itemPackages.MinPrice("USD"),
		MaxPriceUSDFloat: itemPackages.MaxPrice("USD"),
		MinPriceUSDStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MinPrice("USD")))),
		MaxPriceUSDStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MaxPrice("USD")))),

		MinPriceRUBFloat: itemPackages.MinPrice("RUB"),
		MaxPriceRUBFloat: itemPackages.MaxPrice("RUB"),
		MinPriceRUBStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MinPrice("RUB")))),
		MaxPriceRUBStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MaxPrice("RUB")))),

		MinPriceAUDFloat: itemPackages.MinPrice("AUD"),
		MaxPriceAUDFloat: itemPackages.MaxPrice("AUD"),
		MinPriceAUDStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MinPrice("AUD")))),
		MaxPriceAUDStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MaxPrice("AUD")))),

		MinPriceGBPFloat: itemPackages.MinPrice("GBP"),
		MaxPriceGBPFloat: itemPackages.MaxPrice("GBP"),
		MinPriceGBPStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MinPrice("GBP")))),
		MaxPriceGBPStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MaxPrice("GBP")))),

		MinPriceEURFloat: itemPackages.MinPrice("EUR"),
		MaxPriceEURFloat: itemPackages.MaxPrice("EUR"),
		MinPriceEURStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MinPrice("EUR")))),
		MaxPriceEURStr:   fmt.Sprintf("%d", int(math.Ceil(itemPackages.MaxPrice("EUR")))),

		MedPriceBTCFloat: (itemPackages.MinPrice("BTC") + itemPackages.MaxPrice("BTC")) / 2,

		ScoreStr:      humanize.Ftoa(score),
		ScoreFloat:    score,
		GroupPackages: pckgs,
	}

	for _, review := range vi.RatingReviews {
		vi.ViewRatingReviews = append(vi.ViewRatingReviews, review.ViewRatingReview(lang))
	}

	return vi
}

func (items Items) ViewItems(lang string) []ViewItem {
	viewItems := []ViewItem{}
	for _, item := range items {
		viewItem := item.ViewItem(lang)
		viewItems = append(viewItems, viewItem)
	}
	return viewItems
}

func (itms ViewItems) GetItemsPage(pagenumber, pagesize int, sortby string) ([]ViewItem, int) {
	var (
		numberOfPages = int(math.Ceil(float64(len(itms)) / float64(pagesize)))
		pageItems     = []ViewItem{}
	)

	if sortby == "price" {
		sort.Sort(ViewItemsByPrice(itms))
	} else if sortby == "popularity" {
		sort.Sort(ViewItemsByPopularity(itms))
	} else if sortby == "date_added" {
		sort.Sort(ViewItemsByDateAdded(itms))
	} else if sortby == "rating" {
		sort.Sort(ViewItemsByRating(itms))
	} else if sortby == "date_logged_in" {
		sort.Sort(ViewItemsByDateLoggedIn(itms))
	}

	for index, group := range itms {
		if index >= pagenumber*pagesize && index < (pagenumber+1)*pagesize {
			pageItems = append(pageItems, group)
		}
		index += 1
	}

	return pageItems, numberOfPages
}
