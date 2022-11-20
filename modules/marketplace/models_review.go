package marketplace

import (
	"github.com/n0kovo/market_test/modules/util"
	"time"
)

/*
	Models
*/

type RatingReview struct {
	Uuid            string `json:"uuid" gorm:"primary_key"`
	UserUuid        string `json:"user_uuid" sql:"index"`
	ItemUuid        string `json:"item_uuid" sql:"index"`
	TransactionUuid string `json:"transaction_uuid" sql:"index"`
	SellerUuid      string `json:"seller_uuid" sql:"index"`

	ItemScore         int    `json:"item_score"`
	ItemReview        string `json:"item_review"`
	SellerScore       int    `json:"seller_score"`
	SellerReview      string `json:"seller_review"`
	MarketplaceScore  int    `json:"marketplace_score"`
	MarketplaceReview string `json:"marketplace_review"`

	User        User        `json:"-"`
	Seller      User        `json:"-" gorm:"ForeignKey:SellerUuid"`
	Item        Item        `json:"-"`
	Transaction Transaction `json:"-"`

	// ORM timestamps
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ViewRatingReview struct {
	RatingReview
	CreatedAtStr string
	ViewItem     ViewItem
	ViewSeller   ViewSeller
}

/*
	Model Interface Implementation
*/

func (r RatingReview) ViewRatingReview(lang string) ViewRatingReview {
	return ViewRatingReview{
		RatingReview: r,
		ViewItem:     r.Item.ViewItem(lang),
		CreatedAtStr: util.HumanizeTime(*r.CreatedAt, lang),
		ViewSeller:   Seller{&r.Seller}.ViewSeller(lang),
	}
}

func (r RatingReview) Validate() error {
	return nil
}

func (i RatingReview) Remove() error {
	return database.Delete(&i).Error
}

func (itm RatingReview) Save() error {
	err := itm.Validate()
	if err != nil {
		return err
	}
	return itm.SaveToDatabase()
}

func (itm RatingReview) SaveToDatabase() error {
	if existing, _ := FindRatingReviewByUuid(itm.Uuid); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

/*
	Queries
*/

func GetAllRatingReviews() []RatingReview {
	var items []RatingReview
	database.Find(&items)
	return items
}

func FindRatingReviewByUuid(uuid string) (*RatingReview, error) {
	var item RatingReview
	err := database.
		Where(&RatingReview{Uuid: uuid}).
		Preload("Item").
		Preload("Item.Packages").
		Preload("Item.User").
		Preload("User").
		Preload("Seller").
		Preload("Transaction").
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func FindRatingReviewsBySellerUuid(uuid string) ([]RatingReview, error) {
	var items []RatingReview

	err := database.
		Where(&RatingReview{SellerUuid: uuid}).
		Preload("Item").
		Preload("User").
		Preload("Seller").
		Preload("Transaction").
		Find(&items).
		Error

	return items, err
}

func FindRatingReviewByTransactionUuid(uuid string) (*RatingReview, error) {
	var item RatingReview
	err := database.
		Where(&RatingReview{TransactionUuid: uuid}).
		Preload("Item").
		Preload("Item.Packages").
		Preload("User").
		Preload("Seller").
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}
