package marketplace

import (
	"errors"
	"math/rand"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
	"time"
)

type Advertising struct {
	Uuid        string    `json:"uuid" gorm:"primary_key"`
	DateCreated time.Time `json:"created_at"`
	Comment     string    `json:"comment"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Status      bool      `json:"status"`

	CountImpressions        int `json:"count_impressions"`
	CurrentCountImpressions int `json:"Current_count_impressions"`

	VendorUuid string `json:"vendor_uuid" sql:"index"`
	ItemUuid   string `json:"item_uuid" sql:"index"`
	Item       Item   `json:"-"`
}

type Advertisings []Advertising

func (ad Advertising) Validate() error {

	if len(ad.Comment) > 50 || len(ad.Comment) == 0 {
		return errors.New("Text very long. Limit 50 character")
	}

	return nil
}

func addAdvertising(Comment string, Count int, VendorUuid string, ItemUuid string) error {
	ad := Advertising{
		Uuid:                    util.GenerateUuid(),
		DateCreated:             time.Now(),
		Comment:                 Comment,
		DateStart:               time.Now(),
		Status:                  true,
		CountImpressions:        Count,
		CurrentCountImpressions: 0,
		VendorUuid:              VendorUuid,
		ItemUuid:                ItemUuid,
	}
	err := ad.Validate()
	if err != nil {
		return err
	}
	err = ad.Save()
	if err != nil {
		return err
	}

	return nil
}

func FindAllAdvertising() ([]Advertising, error) {
	var ads []Advertising
	err := database.
		Preload("Item").
		Order("date_created ASC").
		Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, err
}

func FindAllActiveAdvertising() ([]Advertising, error) {
	var ads []Advertising
	err := database.
		Where("status = true").
		Preload("Item").
		Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, err
}

func FindAdvertisingByUuid(uuid string) (*Advertising, error) {
	var ad Advertising
	err := database.
		Where("advertisings.item_uuid = ? and advertisings.status = true", uuid).
		Preload("Item").
		Last(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, err
}

func FindAdvertisingByVendor(uuid string) ([]Advertising, error) {
	var ads []Advertising
	err := database.
		Where(&Advertising{VendorUuid: uuid}).
		Preload("Item").
		Find(&ads).Error
	if err != nil {
		return nil, err
	}

	return ads, err
}

func FindAdvertisingByItem(uuid string) (*Advertising, error) {
	var ad Advertising
	err := database.
		Where(&Advertising{ItemUuid: uuid}).
		Preload("Item").
		First(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, err
}

func GetAdvertisings(count int) ([]Advertising, error) {
	var result []Advertising
	ads, err := FindAllActiveAdvertising()
	if err != nil {
		return nil, err
	}

	if len(ads) < count {
		count = len(ads)
	}
	for i := 0; i < count; i++ {
		luckyNumber := rand.Intn(len(ads))
		err := ads[luckyNumber].AddImpressions()
		if err != nil {
			return nil, err
		}
		result = append(result, ads[luckyNumber])
	}

	return result, nil
}

func (ad *Advertising) AddImpressions() error {
	ad.CurrentCountImpressions++
	if ad.CurrentCountImpressions >= ad.CountImpressions {
		ad.Status = false
		ad.DateEnd = time.Now()
		return database.Save(&ad).Error
	}
	return database.Model(&ad).UpdateColumn("CurrentCountImpressions", ad.CurrentCountImpressions).Error
}

func (ad Advertising) Save() error {
	err := ad.Validate()
	if err != nil {
		return err
	}
	return ad.SaveToDatabase()
}

func (ad *Advertising) SaveToDatabase() error {
	if existing, _ := FindAdvertisingByUuid(ad.Uuid); existing == nil {
		return database.Create(ad).Error
	}
	return database.Save(ad).Error
}
