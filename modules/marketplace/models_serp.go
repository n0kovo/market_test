package marketplace

import (
	"fmt"
	"math"
	"sort"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type AvaiableItem struct {
	ItemUuid string `json:"item_uuid"`

	MinPrice float64 `json:"-"`
	MaxPrice float64 `json:"-"`
	Currency string  `json:"-"`

	VendorUuid                      string    `json:"vendor_uuid"`
	VendorUsername                  string    `json:"vendor_username"`
	VendorDescription               string    `json:"vendor_description"`
	VendorLanguage                  string    `json:"vendor_language"`
	Premium                         bool      `json:"vendor_is_premium"`
	PremiumPlus                     bool      `json:"vendor_is_premium_plus"`
	IsTrustedSeller                 bool      `json:"vendor_is_trusted"`
	LastLoginDate                   time.Time `json:"vendor_last_login_date"`
	RegistrationDate                time.Time `json:"vendor_registration_date"`
	BitcoinMultisigPublicKeyEnabled bool      `json:"-"`

	Type                       string    `json:"type"`
	ItemCreatedAt              time.Time `json:"item_created_at"`
	ItemName                   string    `json:"item_name"`
	ItemDescription            string    `json:"item_description"`
	ItemCategoryId             int       `json:"item_category_id"`
	ParentItemCategoryId       int       `json:"item_parent_category_id"`
	ParentParentItemCategoryId int       `json:"item_parent_parent_category_id"`

	SellerScore      float64 `json:"vendor_score"`
	SellerScoreCount int     `json:"vendor_score_count"`
	ItemScore        float64 `json:"item_score"`
	ItemScoreCount   int     `json:"item_score_count"`

	VendorBitcoinTxNumber float64 `json:"-"`
	VendorBitcoinTxVolume float64 `json:"-"`
	ItemBitcoinTxNumber   float64 `json:"-"`
	ItemBitcoinTxVolume   float64 `json:"-"`

	VendorEthereumTxNumber float64 `json:"-"`
	VendorEthereumTxVolume float64 `json:"-"`
	ItemEthereumTxNumber   float64 `json:"-"`
	ItemEthereumTxVolume   float64 `json:"-"`

	CountryNameEnShippingFrom string `json:"country_shipping_from"`
	CountryNameEnShippingTo   string `json:"country_shipping_to"`
	DropCityId                int    `json:"geoname_id"`

	Price map[string][2]float64 `json:"price"`

	GeoCity        City    `gorm:"ForeignKey:DropCityId" json:"-"`
	GeoCountryFrom Country `gorm:"ForeignKey:CountryNameEnShippingFrom" json:"-"`
	GeoCountryTo   Country `gorm:"ForeignKey:CountryNameEnShippingTo" json:"-"`
}

type Vendor struct {
	Username          string
	LastLoginDate     time.Time
	RegistrationDate  time.Time
	Premium           bool
	PremiumPlus       bool
	BitcoinTxNumber   float64
	BitcoinTxVolume   float64
	EthereumTxNumber  float64
	EthereumTxVolume  float64
	VendorScore       float64
	VendorDescription string
	Language          string
	IsTrustedSeller   bool
}

type Vendors []Vendor

/*
	Currency Rates
*/

func (ai AvaiableItem) GetPrice(currency string) [2]float64 {
	return [2]float64{
		ai.MinPrice / GetCurrencyRate(currency, ai.Currency),
		ai.MaxPrice / GetCurrencyRate(currency, ai.Currency),
	}
}

/*
	Util
*/

func isNumberBtwn(number, a, b float64) bool {
	return number >= a && number < b
}

func txVolumeApprox(n float64, currency string) string {
	if n < 0.1 {
		return "< 0.1 " + currency
	} else if isNumberBtwn(n, 0.1, 0.5) {
		return "0.1-0.5 " + currency
	} else if isNumberBtwn(n, 0.5, 1.0) {
		return "0.5-1 " + currency
	} else if isNumberBtwn(n, 1.0, 2.0) {
		return "1-2 " + currency
	} else if isNumberBtwn(n, 2.0, 5.0) {
		return "2-5 " + currency
	} else if isNumberBtwn(n, 5.0, 10.0) {
		return "5-10 " + currency
	} else if n >= 10.0 {
		return "10+ " + currency
	}
	return ""
}

func txNumberApprox(t float64) string {
	if t < 5.0 {
		return "< 5"
	} else if isNumberBtwn(t, 5.0, 10.0) {
		return "5-10"
	} else if isNumberBtwn(t, 10.0, 20.0) {
		return "10-20"
	} else if isNumberBtwn(t, 20.0, 30.0) {
		return "20-30"
	} else if isNumberBtwn(t, 30.0, 40.0) {
		return "30-40"
	} else if isNumberBtwn(t, 40.0, 50.0) {
		return "40-50"
	} else if isNumberBtwn(t, 50.0, 100.0) {
		return "50-100"
	} else if t >= 100.0 {
		return "100+"
	}
	return ""
}

/*
	View Item
*/

type ViewAvailableItem struct {
	*AvaiableItem
	IsOnline                  bool      `json:"vendor_is_online"`
	LastLoginDateStr          string    `json:"vendor_last_login_date"`
	RegistrationDateStr       string    `json:"vendor_registration_date"`
	PriceRangeStr             [2]string `json:"price_range"`
	PriceStr                  string    `json:"price"`
	VendorBitcoinTxNumberStr  string    `json:"vendor_btc_tx_number"`
	VendorBitcoinTxVolumeStr  string    `json:"vendor_btc_tx_volume"`
	ItemBitcoinTxNumberStr    string    `json:"item_btc_tx_number"`
	ItemBitcoinTxVolumeStr    string    `json:"item_btc_tx_volume"`
	VendorEthereumTxNumberStr string    `json:"vendor_eth_tx_number"`
	VendorEthereumTxVolumeStr string    `json:"item_btc_tx_volume"`
	ItemEthereumTxNumberStr   string    `json:"item_eth_tx_number"`
	ItemEthereumTxVolumeStr   string    `json:"item_btc_tx_volume"`
}

func (ai AvaiableItem) ViewAvailableItem(lang, currency string) ViewAvailableItem {

	price := ai.Price[currency]

	vai := ViewAvailableItem{
		AvaiableItem:              &ai,
		ItemBitcoinTxNumberStr:    txNumberApprox(ai.ItemBitcoinTxNumber),
		ItemBitcoinTxVolumeStr:    txVolumeApprox(ai.ItemBitcoinTxVolume, "BTC"),
		ItemEthereumTxNumberStr:   txNumberApprox(ai.ItemEthereumTxNumber),
		ItemEthereumTxVolumeStr:   txVolumeApprox(ai.ItemEthereumTxVolume, "ETH"),
		LastLoginDateStr:          util.HumanizeTime(ai.LastLoginDate, lang),
		PriceRangeStr:             [2]string{fmt.Sprintf("%d", int(math.Ceil(price[0]))), fmt.Sprintf("%d", int(math.Ceil(price[1])))},
		RegistrationDateStr:       util.HumanizeTime(ai.RegistrationDate, lang),
		VendorBitcoinTxNumberStr:  txNumberApprox(ai.VendorBitcoinTxNumber),
		VendorBitcoinTxVolumeStr:  txVolumeApprox(ai.VendorBitcoinTxVolume, "BTC"),
		VendorEthereumTxNumberStr: txNumberApprox(ai.VendorEthereumTxNumber),
		VendorEthereumTxVolumeStr: txVolumeApprox(ai.VendorEthereumTxVolume, "ETH"),
	}

	if currency == "BTC" || currency == "ETH" || currency == "BCH" {
		vai.PriceRangeStr = [2]string{
			fmt.Sprintf("%f", price[0]),
			fmt.Sprintf("%f", price[1]),
		}
	}

	if price[0] == price[1] {
		vai.PriceStr = vai.PriceRangeStr[0]
	}

	return vai
}

type ViewVendor struct {
	*Vendor
	IsOnline            bool
	BitcoinTxVolumeStr  string
	BitcoinTxNumberStr  string
	EthereumTxVolumeStr string
	EthereumTxNumberStr string
	VendorScoreStr      string
	LastLoginDateStr    string
	RegistrationDateStr string
}

func (v Vendor) ViewVendor(lang string) ViewVendor {
	return ViewVendor{
		BitcoinTxNumberStr:  txNumberApprox(v.BitcoinTxNumber),
		BitcoinTxVolumeStr:  txVolumeApprox(v.BitcoinTxVolume, "BTC"),
		EthereumTxNumberStr: txNumberApprox(v.EthereumTxNumber),
		EthereumTxVolumeStr: txVolumeApprox(v.EthereumTxVolume, "ETH"),
		LastLoginDateStr:    util.HumanizeTime(v.LastLoginDate, lang),
		RegistrationDateStr: util.HumanizeTime(v.RegistrationDate, lang),
		Vendor:              &v,
	}
}

type AvailableItems []AvaiableItem

func (ais AvailableItems) ViewAvailableItems(lang, currency string) []ViewAvailableItem {
	var vais []ViewAvailableItem
	for _, ai := range ais {
		vai := ai.ViewAvailableItem(lang, currency)
		vais = append(vais, vai)
	}
	return vais
}

func (vs Vendors) ViewVendors(lang string) []ViewVendor {
	var vvs []ViewVendor
	for _, v := range vs {
		vv := v.ViewVendor(lang)
		vvs = append(vvs, vv)
	}
	return vvs
}

/*
	Collection Fields
*/

func (ais AvailableItems) Sort(sortyBy string) AvailableItems {

	var sortByFunc func(int, int) bool

	switch sortyBy {
	case "date_logged_in":
		sortByFunc = func(i, j int) bool {
			return ais[i].LastLoginDate.After(ais[j].LastLoginDate)
		}
	case "price":
		sortByFunc = func(i, j int) bool {
			if ais[i].PremiumPlus != ais[j].PremiumPlus { // by premium status
				return ais[i].PremiumPlus
			} else if ais[i].Premium != ais[j].Premium { // by premium status
				return ais[i].Premium
			} else if ais[i].ItemScore != ais[j].ItemScore { // by score
				return ais[i].ItemScore > ais[j].ItemScore
			} else { // by price
				return ais[i].MinPrice < ais[j].MaxPrice
			}
		}
	case "popularity":
		sortByFunc = func(i, j int) bool {
			if ais[i].Premium != ais[j].Premium { // by premium status
				return ais[i].Premium
			} else { // by price
				return ais[i].ItemBitcoinTxNumber+ais[i].ItemEthereumTxNumber > ais[j].ItemBitcoinTxNumber+ais[j].ItemEthereumTxNumber
			}
		}
	case "date_added":
		sortByFunc = func(i, j int) bool {
			return ais[i].ItemCreatedAt.Before(ais[j].ItemCreatedAt)
		}
	case "rating":
		sortByFunc = func(i, j int) bool {
			return ais[i].ItemScore*float64(ais[i].ItemBitcoinTxNumber+ais[i].ItemBitcoinTxNumber) >
				ais[j].ItemScore*float64(ais[j].ItemBitcoinTxNumber+ais[j].ItemBitcoinTxNumber)

		}
	default:
		sortByFunc = func(i, j int) bool { return true }
	}

	sort.Slice(ais, sortByFunc)
	return ais
}

func (ais AvailableItems) Where(predicate func(AvaiableItem) bool) AvailableItems {
	newAis := AvailableItems{}
	for i, _ := range ais {
		if predicate(ais[i]) {
			newAis = append(newAis, ais[i])
		}
	}
	return newAis
}

func (ais AvailableItems) Filter(category, dropCityId int,
	packageType, query, to, from, accountType string) AvailableItems {

	var searchResults []string
	if query != "" {
		searchResults = SearchItems(query)
	}

	categoryPredicate := func(ai AvaiableItem) bool {
		if category != 0 {
			return ai.ItemCategoryId == category || ai.ParentItemCategoryId == category || ai.ParentParentItemCategoryId == category
		}
		return true
	}

	accountTypePredicate := func(ai AvaiableItem) bool {
		if accountType == "premium" {
			return ai.Premium
		}
		if accountType == "premium_plus" {
			return ai.PremiumPlus
		}
		return true
	}

	typePredicate := func(ai AvaiableItem) bool {
		if packageType != "" && packageType != "all" {
			return ai.Type == packageType
		}
		return true
	}

	queryPredicate := func(ai AvaiableItem) bool {
		if query != "" && !inSet(ai.ItemUuid, searchResults) {
			return false
		}
		return true
	}

	shipingToPredicate := func(ai AvaiableItem) bool {
		if to == "" {
			return true
		}
		return ai.CountryNameEnShippingTo == to
	}

	shipingFromPredicate := func(ai AvaiableItem) bool {
		if from == "" {
			return true
		}
		return ai.CountryNameEnShippingFrom == from
	}

	cityPredicate := func(ai AvaiableItem) bool {
		if dropCityId == 0 {
			return true
		}
		return ai.DropCityId == dropCityId
	}

	filteredAvailableItems := ais.
		Where(typePredicate).
		Where(categoryPredicate).
		Where(queryPredicate).
		Where(shipingToPredicate).
		Where(shipingFromPredicate).
		Where(cityPredicate).
		Where(accountTypePredicate)

	return filteredAvailableItems
}

func (ais AvailableItems) DropCitiesList() []City {
	locationMap := map[int]City{}
	for _, a := range ais {
		locationMap[a.DropCityId] = a.GeoCity
	}
	locations := []City{}
	for _, city := range locationMap {
		locations = append(locations, city)
	}
	return locations
}

func (ais AvailableItems) ShippingToList() []string {
	locationMap := map[string]bool{}
	for _, a := range ais {
		locationMap[a.CountryNameEnShippingTo] = true
	}
	locations := []string{}
	for l, _ := range locationMap {
		locations = append(locations, l)
	}
	sort.Strings(locations)
	return locations
}

func (ais AvailableItems) ShippingFromList() []string {
	locationMap := map[string]bool{}
	for _, a := range ais {
		locationMap[a.CountryNameEnShippingFrom] = true
	}
	locations := []string{}
	for l, _ := range locationMap {
		locations = append(locations, l)
	}
	sort.Strings(locations)
	return locations
}

func (ais AvailableItems) VendorList() Vendors {
	vendorsMap := map[string]Vendor{}
	for _, a := range ais {
		vendorsMap[a.VendorUsername] = Vendor{
			Username:          a.VendorUsername,
			LastLoginDate:     a.LastLoginDate,
			RegistrationDate:  a.RegistrationDate,
			Premium:           a.Premium,
			PremiumPlus:       a.PremiumPlus,
			BitcoinTxNumber:   a.VendorBitcoinTxNumber,
			BitcoinTxVolume:   a.VendorBitcoinTxVolume,
			EthereumTxNumber:  a.VendorEthereumTxNumber,
			EthereumTxVolume:  a.VendorEthereumTxVolume,
			VendorScore:       a.SellerScore,
			VendorDescription: a.VendorDescription,
			Language:          a.VendorLanguage,
			IsTrustedSeller:   a.IsTrustedSeller,
		}
	}

	vendors := []Vendor{}
	for _, v := range vendorsMap {
		vendors = append(vendors, v)
	}

	return Vendors(vendors)
}

func bool2int(b bool) float64 {
	if b {
		return 1.0
	} else {
		return 0.0
	}
}

func (vvs Vendors) Sort(sortyBy string) Vendors {

	var sortByFunc func(int, int) bool

	switch sortyBy {
	case "date_logged_in":
		sortByFunc = func(i, j int) bool {
			return vvs[i].LastLoginDate.After(vvs[j].LastLoginDate)
		}
	case "popularity":
		sortByFunc = func(i, j int) bool {
			return (vvs[i].BitcoinTxNumber+vvs[i].EthereumTxNumber)*(1+bool2int(vvs[i].IsTrustedSeller)) >
				(vvs[j].BitcoinTxNumber+vvs[j].EthereumTxNumber)*(1+bool2int(vvs[j].IsTrustedSeller))
		}
	case "date_added":
		sortByFunc = func(i, j int) bool {
			return vvs[i].RegistrationDate.After(vvs[j].RegistrationDate)
		}
	case "rating":
		sortByFunc = func(i, j int) bool {
			if vvs[i].IsTrustedSeller != vvs[j].IsTrustedSeller {
				return vvs[i].IsTrustedSeller
			} else if vvs[i].PremiumPlus != vvs[j].PremiumPlus { // by premium status
				return vvs[i].PremiumPlus
			} else if vvs[i].Premium != vvs[j].Premium { // by premium status
				return vvs[i].Premium
			} else {
				return vvs[i].VendorScore > vvs[j].VendorScore
			}
		}
	default:
		sortByFunc = func(i, j int) bool { return true }
	}

	sort.Slice(vvs, sortByFunc)
	return vvs
}

/*
	Database Queries
*/

func FindAvailableItems(
	deliveryType, accountType, vendorUuid string,
	categoryId, page, pageSize int,
) AvailableItems {
	items := []AvaiableItem{}

	query := database.Table("v_available_items")

	if deliveryType != "" {
		query = query.Where("type=?", vendorUuid)
	}

	if vendorUuid != "" {
		query = query.Where("vendor_uuid=?", vendorUuid)
	}

	if categoryId != 0 {
		query = query.Where("item_category_id=?", categoryId)
	}
	if accountType != "" {
		if accountType == "premium_plus" {
			query = query.Where("(premium=true OR premium_plus=true)")
		}
		if accountType == "premium" {
			query = query.Where("(premium=true)")
		}
	}

	query = query.Preload("GeoCity")

	if page != 0 && pageSize != 0 {
		query = query.
			Offset(pageSize * page).
			Limit(pageSize)
	}

	query.Find(&items)

	for i, _ := range items {
		items[i].Price = map[string][2]float64{
			"AUD": items[i].GetPrice("AUD"),
			"BTC": items[i].GetPrice("BTC"),
			"ETH": items[i].GetPrice("ETH"),
			"EUR": items[i].GetPrice("EUR"),
			"GBP": items[i].GetPrice("GBP"),
			"RUB": items[i].GetPrice("RUB"),
			"USD": items[i].GetPrice("USD"),
		}

		items[i].ItemScore = float64(int(items[i].ItemScore*100)) / float64(100.0)
		items[i].SellerScore = float64(int(items[i].SellerScore*100)) / float64(100.0)
	}

	return AvailableItems(items)
}

func CountAvailableItems(
	deliveryType, accountType string,
	categoryId int,
) int {
	count := 0

	query := database.Table("v_available_items")

	if deliveryType != "" {
		query = query.Where("type=?", accountType)
	}

	if categoryId != 0 {
		query = query.Where("item_category_id=?", categoryId)
	}

	if accountType != "" {
		if accountType == "premium_plus" {
			query = query.Where("(premium=true OR premium_plus=true)")
		}
		if accountType == "premium" {
			query = query.Where("(premium=true)")
		}
	}

	query.Count(&count)

	return count
}

/*
	In-Memory Queries
*/

func CacheGetAvailableItems() AvailableItems {

	key := "available-items"
	cAvailableItems, _ := gc.Get(key)
	if cAvailableItems == nil {
		availableItems := FindAvailableItems("", "", "", 0, 0, 0)
		gc.Set(key, availableItems)
		return availableItems
	}

	return cAvailableItems.(AvailableItems)
}

/*
	Utils
*/

func inSet(s string, ss []string) bool {
	for _, i := range ss {
		if i == s {
			return true
		}
	}
	return false
}

/*
	Database Views
*/

func setupAvailableItemsView() {

	database.Exec("DROP VIEW IF EXISTS v_available_items CASCADE;")
	database.Exec(`
		CREATE VIEW v_available_items AS (
			select 
				v_packages.item_uuid,
				v_packages.drop_city_id,
				v_packages.country_name_en_shipping_from,
				v_packages.country_name_en_shipping_to,
				v_packages.currency,
				min(v_packages.price) as min_price,
				max(v_packages.price) as max_price,
				users.uuid as vendor_uuid,
				users.description as vendor_description,
				users.username as vendor_username,
				users.language as vendor_language,
				users.premium,
				users.premium_plus,
				users.is_trusted_seller,
				users.last_login_date,
				users.registration_date,
				users.bitcoin_multisig_public_key != '' as bitcoin_multisig_public_key_enabled, 
				type,
				items.created_at as item_created_at, 
				items.name as item_name,
				items.description as item_description,
				items.item_category_id,
				ic_parent.id as parent_item_category_id,
				ic_parent.parent_id as parent_parent_item_category_id,
				COALESCE(avg(r1.seller_score), 0) as seller_score,
				COALESCE(count(r1.seller_score), 0) as seller_score_count,
				COALESCE(avg(r2.item_score), 0) as item_score,
				COALESCE(count(r2.item_score), 0) as item_score_count,
				AVG(COALESCE(v_vendor_bitcoin_tx_stats.tx_number, 0)) as vendor_bitcoin_tx_number, 
				AVG(COALESCE(v_vendor_bitcoin_tx_stats.tx_volume, 0)) as vendor_bitcoin_tx_volume,
				AVG(COALESCE(v_vendor_ethereum_tx_stats.tx_number, 0)) as vendor_ethereum_tx_number, 
				AVG(COALESCE(v_vendor_ethereum_tx_stats.tx_volume, 0)) as vendor_ethereum_tx_volume,
				AVG(COALESCE(v_item_bitcoin_tx_stats.tx_number, 0)) as item_bitcoin_tx_number,
				AVG(COALESCE(v_item_bitcoin_tx_stats.tx_volume, 0)) as item_bitcoin_tx_volume,
				AVG(COALESCE(v_item_ethereum_tx_stats.tx_number, 0)) as item_ethereum_tx_number,
				AVG(COALESCE(v_item_ethereum_tx_stats.tx_volume, 0)) as item_ethereum_tx_volume
			from v_packages 
			join items on items.uuid = v_packages.item_uuid
			left join v_vendor_bitcoin_tx_stats on v_vendor_bitcoin_tx_stats.seller_uuid = items.user_uuid
			left join v_vendor_ethereum_tx_stats on v_vendor_ethereum_tx_stats.seller_uuid = items.user_uuid
			left join v_item_bitcoin_tx_stats on v_item_bitcoin_tx_stats.item_uuid = items.uuid
			left join v_item_ethereum_tx_stats on v_item_ethereum_tx_stats.item_uuid = items.uuid
			join users on users.uuid = items.user_uuid
			left join item_categories ic on ic.id = items.item_category_id
			left join item_categories ic_parent on ic_parent.id = ic.parent_id
			left join rating_reviews r1 on r1.seller_uuid = items.user_uuid
			left join rating_reviews r2 on r2.item_uuid = v_packages.item_uuid
			WHERE items.deleted_at IS NULL AND users.banned=false and v_packages.deleted_at IS NULL
			group by 
				users.uuid,
				users.username, 
				users.premium, 
				users.premium_plus, 
				users.is_trusted_seller,
				users.last_login_date, 
				users.registration_date, 
				users.bitcoin_multisig_public_key,
				v_packages.item_uuid, 
				v_packages.drop_city_id,
				v_packages.country_name_en_shipping_from,
				v_packages.country_name_en_shipping_to,
				v_packages.currency,
				r1.seller_uuid, 
				r2.item_uuid, 
				type, 
				items.created_at,
				items.name, 
				items.description,
				items.item_category_id,
				parent_item_category_id,
				parent_parent_item_category_id
	);`)
}
