package marketplace

import (
	"encoding/json"
	"errors"
	"strconv"
)

type City struct {
	ID           int    `json:"geonameid" gorm:"primary_key"`
	NameEn       string `json:"name" sql:"index"`
	CountryEn    string `json:"country" sql:"index"`
	SubCountryEn string `json:"subcountry" sql:"index"`

	Country Country `gorm:"AssociationForeignKey:NameEn" json:"-"`
}

type Country struct {
	NameEn string `json:"name_en" gorm:"primary_key"`
}

/*
	Model Interface Implementation
*/

func (r City) Validate() error {
	if r.ID == 0 {
		return errors.New("Empty ID")
	}
	if r.NameEn == "" {
		return errors.New("NameEn must be nonempy")
	}
	if r.CountryEn == "" {
		return errors.New("CountryEn must be nonempy")
	}
	if r.SubCountryEn == "" {
		return errors.New("SubCountryEn must be nonempy")
	}
	return nil
}

func (i City) Remove() error {
	return database.Delete(&i).Error
}

func (itm City) Save() error {
	err := itm.Validate()
	if err != nil {
		return err
	}
	return itm.SaveToDatabase()
}

func (itm City) SaveToDatabase() error {
	if existing, _ := FindCityById(itm.ID); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

func (r Country) Validate() error {
	if r.NameEn == "" {
		return errors.New("Invalid Name EN")
	}
	return nil
}

func (i Country) Remove() error {
	return database.Delete(&i).Error
}

func (itm Country) Save() error {
	err := itm.Validate()
	if err != nil {
		return err
	}
	return itm.SaveToDatabase()
}

func (itm Country) SaveToDatabase() error {
	if existing, _ := FindCountryByNameEn(itm.NameEn); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

/*
	JSON
*/

func (u *City) UnmarshalJSON(data []byte) error {

	ct := struct {
		ID           string `json:"geonameid"`
		NameEn       string `json:"name" gorm:"primary_key"`
		CountryEn    string `json:"country" sql:"index"`
		SubCountryEn string `json:"subcountry" sql:"index"`
	}{}

	if err := json.Unmarshal(data, &ct); err != nil {
		return err
	}

	u.NameEn = ct.NameEn
	u.CountryEn = ct.CountryEn
	u.SubCountryEn = ct.SubCountryEn

	i, err := strconv.ParseInt(ct.ID, 10, 64)
	if err != nil {
		return err
	}
	u.ID = int(i)

	return err
}

/*
	Queries
*/

func GetAllCities() []City {
	cities := []City{}
	database.Find(&cities)
	return cities
}

func FindCityById(id int) (*City, error) {
	var item City
	err := database.
		First(&item, "id=?", id).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func FindCitiesByCountryNameEn(name string) []City {
	var items []City
	database.
		Where("country_en=?", name).
		Order("name_en").
		Find(&items)
	return items
}

func FindCountryByNameEn(name string) (*Country, error) {
	var item Country
	err := database.
		First(&item, "name_en=?", name).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func GetAllEnglishCountryNamesFromCities() []Country {
	var countries []Country

	database.
		Raw("SELECT DISTINCT(country_en) as name_en from cities").
		Scan(&countries)

	return countries
}

func GetAllCountries() []Country {
	countries := []Country{}
	database.Order("name_en").Find(&countries)
	return countries
}
