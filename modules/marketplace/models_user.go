package marketplace

import (
	"errors"
	"html/template"
	"math"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jinzhu/gorm"
	"github.com/o1egl/govatar"
	"github.com/russross/blackfriday"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

var (
	onlineDuration, _ = time.ParseDuration(MARKETPLACE_SETTINGS.OnlineDuration)
)

/*
	Models
*/

type User struct {
	// Base shit
	Uuid           string `json:"uuid" gorm:"primary_key"`
	Username       string `json:"username" sql:"type:varchar(16);unique_index"`
	PassphraseHash string `json:"passphrase_hash_v1"`

	// Login dates
	RegistrationDate time.Time  `json:"registration_date" gorm:"index"`
	LastLoginDate    *time.Time `json:"last_login_date" gorm:"index"`

	// Settings
	Language string `json:"language" gorm:"index"`
	Currency string `json:"currency" gorm:"index"`

	// Profile
	Bitcoin                  string `json:"btc_address"`
	BitcoinMultisigPublicKey string `json:"btc_multisig_public_key"`
	Ethereum                 string `json:"ethereum_address"`
	BitcoinCash              string `json:"bch_address"`

	// Contacts
	Bitmessage      string `json:"bitmessage"`
	Tox             string `json:"tox"`
	Email           string `json:"email"`
	Pgp             string `sql:"type:varchar(8192);"`
	Description     string `json:"description" sql:"size:140"`
	LongDescription string `json:"long_description" sql:"size:2048"`
	InviteCode      string `json:"invite_code" sql:"unique_index"`

	// Misc settings
	TwoFactorAuthentication bool `json:"2fa_enabled" gorm:"index"`
	Premium                 bool `json:"is_premium" gorm:"index"`
	PremiumPlus             bool `json:"is_premium_plus" gorm:"index"`
	HasTopBanner            bool `json:"has_top_banner" gorm:"index"`
	Banned                  bool `json:"is_banned" gorm:"index"`
	PossibleScammer         bool `json:"is_possible_scammer" gorm:"index"`
	VacationMode            bool `json:"vacation_mode" gorm:"index"`
	HasAvatar               bool `json:"has_avatar" gorm:"index"`

	// Relations
	Items                  []Item                  `json:"-"`
	RatingReviews          []RatingReview          `json:"-"`
	ShippingOptions        []ShippingOption        `json:"-"`
	UserBitcoinWallets     []UserBitcoinWallet     `json:"-"`
	UserEthereumWallets    []UserEthereumWallet    `json:"-"`
	UserBitcoinCashWallets []UserBitcoinCashWallet `json:"-"`

	// Roles
	IsSeller        bool `json:"is_seller" gorm:"index"`
	IsTrustedSeller bool `json:"is_trustedseller" gorm:"index"`
	IsTester        bool `json:"is_tester" gorm:"index"`
	IsModerator     bool `json:"is_moderator" gorm:"index"`
	IsAdmin         bool `json:"is_admin" gorm:"index"`
	IsStaff         bool `json:"is_staff" gorm:"index"`

	// Other Bools
	HasRequestedVerification bool `json:"has_requested_verification" gorm:"index"`

	// Integrations
	MattermostUsername string `json:"mattermost_username"`

	// Index
	InviterUuid string `json:"inviter_uuid" gorm:"index"`

	// Support User UUID
	SupporterUuid string `json:"supporter_uuid" gorm:"index"`
	TrusteeUuid   string `json:"trustee_uuid" gorm:"index"`

	// ORM timestamps
	CreatedAt *time.Time `json:"crated_at" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type ExtendedUser struct {
	User
	InviterUsername               string
	InviterCount                  string
	BitcoinBalance                float64
	BitcoinUnconfirmedBalance     float64
	BitcoinCashBalance            float64
	BitcoinCashUnconfirmedBalance float64
	EthereumBalance               float64
	NumberOfSupportMessages       int
	LastMessageByStaff            bool
	SupportUserUsername           string
}

type Seller struct {
	*User
}

type Users []User

type ExtendedUsers []ExtendedUser

type Sellers []Seller

type UserSettingsHistory struct {
	gorm.Model

	UserUuid string    `json:"user_uuid"`
	Datetime time.Time `json:"datetime"`
	Action   string    `json:"action" sql:"size:1024"`
	Type     string    `json:"string"`
	User     User
}

/*
	Model Methods
*/

func (u User) CheckPassphrase(passphrase string) bool {
	return u.PassphraseHash == util.PasswordHashV1(u.Username, passphrase)
}

func (u User) String() string {
	return u.Username
}

func (u User) NumberOfInvitedUsers() int {
	var count int
	database.
		Model(&User{}).
		Where(&User{InviterUuid: u.Uuid}).
		Count(&count)
	return count
}

// func (u User) Balance() float64

func (u User) Validate() error {
	if !usernameRegexp.MatchString(u.Username) {
		return errors.New("Username is not valid")
	}
	if u.Bitmessage != "" && !bitmessageRegexp.MatchString(u.Bitmessage) {
		return errors.New("Invalid Bitmessage")
	}
	if u.Pgp != "" {
		fromBlock, err := armor.Decode(strings.NewReader(u.Pgp))
		if err != nil || fromBlock.Type != openpgp.PublicKeyType {
			return errors.New("Invalid PGP Key")
		}
	}
	if u.Bitcoin != "" && !bitcoinRegexp.MatchString(u.Bitcoin) {
		return errors.New("Ivalid BTC address")
	}
	if u.Ethereum != "" && !ethereumRegexp.MatchString(u.Ethereum) {
		return errors.New("Ivalid ETH address")
	}
	if u.TwoFactorAuthentication && u.Pgp == "" {
		return errors.New("Fill PGP for 2FA")
	}
	if u.Email != "" && !emailRegexp.MatchString(u.Email) {
		return errors.New("Invalid email/xmpp address")
	}
	if u.InviteCode == "" {
		u.InviteCode = util.GenerateUuid()
	}
	return nil
}

func (s Seller) Score() float64 {
	score := float64(0.0)
	for _, r := range s.RatingReviews {
		score += float64(r.SellerScore)
	}
	if len(s.RatingReviews) > 0 {
		score /= float64(len(s.RatingReviews))
	}
	return math.Ceil(score*100) / float64(100.0)
}

/*
	Relations
*/
func (u User) Iniviter() *User {
	if u.InviterUuid == "" {
		return nil
	}
	user, _ := FindUserByUuid(u.InviterUuid, false)
	return user
}

func (u User) Fingerprint() string {
	return util.Fingerprint(u.Pgp)
}

func (s Seller) FindItemsForSeller() Items {
	return FindItemsForSeller(s.Uuid)
}

func (u User) FindUserBitcoinWallets() UserBitcoinWallets {
	var uw UserBitcoinWallets

	database.
		Where(&UserBitcoinWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateBitcoinWallet(u)
		if err != nil {
			return uw
		}

		return UserBitcoinWallets{
			*userWallet,
		}
	}

	return UserBitcoinWallets(uw)
}

func (u User) FindUserBitcoinCashWallets() UserBitcoinCashWallets {
	var uw UserBitcoinCashWallets

	database.
		Where(&UserBitcoinCashWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateBitcoinCashWallet(u)
		if err != nil {
			return uw
		}

		return UserBitcoinCashWallets{
			*userWallet,
		}
	}

	return UserBitcoinCashWallets(uw)
}

func (u User) FindRecentBitcoinWallet() UserBitcoinWallet {
	return u.FindUserBitcoinWallets()[0]
}

func (u User) FindRecentBitcoinCashWallet() UserBitcoinCashWallet {
	return u.FindUserBitcoinCashWallets()[0]
}

func (u User) FindUserEthereumWallets() UserEthereumWallets {
	var uw UserEthereumWallets

	database.
		Where(&UserEthereumWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateEthereumWallet(u)
		if err != nil {
			return uw
		}
		return UserEthereumWallets{*userWallet}
	}

	return UserEthereumWallets(uw)
}

func (u User) FindRecentEthereumWallet() UserEthereumWallet {
	return u.FindUserEthereumWallets()[0]
}

/*
	Queries
*/

func SettingsChangeHistoryByUser(uuid string) []UserSettingsHistory {
	var history []UserSettingsHistory
	database.
		Where(&UserSettingsHistory{UserUuid: uuid}).
		Find(&history)
	return history
}

func CreateUser(username string, passphrase string) (*User, error) {
	user, _ := FindUserByUsername(username)
	if user != nil {
		return nil, errors.New("Username is not unique")
	}
	invitation, _ := FindInvitationByUsername(username)
	if invitation != nil {
		return nil, errors.New("Username is reserved")
	}
	user = &User{
		Uuid:             util.GenerateUuid(),
		InviteCode:       util.GenerateUuid(),
		PassphraseHash:   util.PasswordHashV1(username, passphrase),
		Username:         username,
		RegistrationDate: time.Now(),
		Currency:         "USD",
	}
	validationError := user.Validate()
	if validationError != nil {
		return nil, validationError
	}
	user.Save()

	return user, nil
}

func CreateUserWithInvitation(invitation Invitation, passphrase string) (*User, error) {
	user, _ := FindUserByUsername(invitation.Username)
	if user != nil {
		return nil, errors.New("Username is not unique")
	}
	user = &User{
		Uuid:             util.GenerateUuid(),
		InviteCode:       util.GenerateUuid(),
		PassphraseHash:   util.PasswordHashV1(invitation.Username, passphrase),
		Username:         invitation.Username,
		InviterUuid:      invitation.InviterUuid,
		RegistrationDate: time.Now(),
		Premium:          true,
	}
	validationError := user.Validate()
	if validationError != nil {
		return nil, validationError
	}
	user.Save()
	invitation.IsActivated = true
	err := invitation.Save()

	return user, err
}

func GetExtendedUsersPage(page, pageSize int, orderBy, query string) []ExtendedUser {

	if orderBy == "last_login" {
		orderBy = "last_login_date desc nulls last"
	} else if orderBy == "invited_users" {
		orderBy = "inviter_count desc nulls last"
	} else if orderBy == "balance" {
		orderBy = "bitcoin_balance desc nulls last, bitcoin_cash_balance desc nulls last, ethereum_balance desc nulls last"
	} else {
		orderBy = "registration_date desc nulls last"
	}

	var users []ExtendedUser
	qry := database.
		Table("v_users").
		Model(&ExtendedUser{}).
		Order(orderBy).
		Offset(page * pageSize).
		Limit(pageSize)

	if query != "" {
		qry = qry.Where("username LIKE ?", "%"+query+"%")
	}

	qry.Find(&users)

	return users
}

func GetInvitedUserPage(inviterUuid string, page, pageSize int) []ExtendedUser {
	orderBy := "registration_date desc nulls last"
	var users []ExtendedUser
	qry := database.
		Table("v_users").
		Where("inviter_uuid = ?", inviterUuid).
		Model(&ExtendedUser{}).
		Order(orderBy).
		Offset(page * pageSize).
		Limit(pageSize)
	qry.Find(&users)
	return users
}

func CountUsers(dt *time.Time) int {
	var count int
	q := database.Table("users")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountVendors(dt *time.Time) int {
	var count int
	q := database.Table("users").Where("is_seller=?", true)
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountVendorsFree(dt *time.Time) int {
	var count int
	q := database.Table("users").Where("is_seller=? and premium=? and premium_plus=?", true, false, false)
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountVendorsPremium(dt *time.Time) int {
	var count int
	q := database.Table("users").Where("is_seller=? and premium=? and premium_plus=?", true, true, false)
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountVendorsPremiumPlus(dt *time.Time) int {
	var count int
	q := database.Table("users").Where("is_seller=? and premium_plus=?", true, true)
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountUsersRegistredAfter(dt time.Time) int {
	var count int
	database.Table("users").Where("registration_date > ?", dt).Count(&count)
	return count
}

func CountUsersActiveAfter(dt time.Time) int {
	var count int
	database.Table("users").Where("last_login_date > ?", dt).Count(&count)
	return count
}

func FindUserByUuid(uuid string, preloadShipping bool) (*User, error) {
	var user User
	q := database.Table("users")

	if preloadShipping {
		q = q.Preload("ShippingOptions")
	}

	err := q.
		First(&user, "uuid = ?", uuid).
		Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	err := database.
		First(&user, "username = ?", username).
		Preload("ShippingOptions").
		Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func FindSellerByUsername(username string) (*Seller, error) {
	var user User
	err := database.
		First(&user, "username = ?", username).
		Preload("Items").
		Preload("ShippingOptions").
		Error
	if err != nil {
		return nil, err
	}
	return &Seller{&user}, err
}

func FindVendors(isPremium, isPremiumPlus, isTrusted *bool) ([]User, error) {
	var users []User
	query := database.Table("users").Where("is_seller = ?", true)

	if isPremium != nil {
		query = query.Where("premium = ?", isPremium)
	}

	if isPremiumPlus != nil {
		query = query.Where("premium_plus = ?", isPremiumPlus)
	}

	if isTrusted != nil {
		query = query.Where("is_trusted_seller = ?", isTrusted)
	}

	query = query.Order("has_requested_verification DESC NULLS LAST, registration_date DESC")

	err := query.Find(&users).Error
	if err != nil {
		return []User{}, err
	}

	return users, err
}

func FindUserByInviteCode(code string) (*User, error) {
	var user User
	err := database.
		First(&user, "invite_code = ?", code).
		Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

/*
	Staff
*/

func FindStaffMemebers() ([]ExtendedUser, error) {
	var users []ExtendedUser
	err := database.
		Table("v_users").
		Where("is_staff=true OR is_admin=true").
		Order("is_admin DESC, is_staff DESC, last_login_date ASC").
		Find(&users).
		Error
	return users, err
}

func FindUncontactedUsers(page, pageSize int) ([]ExtendedUser, error) {
	var users []ExtendedUser
	err := database.
		Table("v_users").
		Where("number_of_support_messages > 1 and last_message_by_staff = false").
		Order("last_message_datetime DESC").
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&users).
		Error
	return users, err
}

func CountUncontactedUsers() int {
	var count int
	database.
		Table("v_users").
		Where("number_of_support_messages > 1 and last_message_by_staff = false").
		Count(&count)
	return count
}

func FindUsersContactedByStaff(page, pageSize int) ([]ExtendedUser, error) {
	var users []ExtendedUser
	err := database.
		Table("v_users").
		Where("last_message_by_staff = true").
		Order("last_message_datetime DESC").
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&users).
		Error
	return users, err
}

func CountUsersContactedByStaff(staffUuid string) int {
	var count int
	database.
		Table("v_users").
		Where("last_message_by_staff = true").
		Order("last_message_datetime DESC").
		Count(&count)
	return count
}

func FindUsersContactedByStaffNeedAnswer(staffUuid string, page, pageSize int) ([]ExtendedUser, error) {
	var users []ExtendedUser
	err := database.
		Table("v_users").
		Where("number_of_support_messages > 1 and last_message_by_staff=false").
		Order("registration_date DESC").
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&users).
		Error
	return users, err
}

func CountUsersContactedByStaffNeedAnswer(staffUuid string) int {
	var count int
	database.
		Table("v_users").
		Where("number_of_support_messages > 1 and last_message_by_staff=false").
		Order("registration_date DESC").
		Count(&count)
	return count
}

/*
	Model Methods
*/

func (u User) Save() error {
	err := u.Validate()
	if err != nil {
		return err
	}
	return u.SaveToDatabase()
}

func (u User) SaveToDatabase() error {
	if existing, _ := FindUserByUuid(u.Uuid, false); existing == nil {
		return database.Create(&u).Error
	}
	return database.Save(&u).Error
}

func (u User) Remove() error {
	return database.Delete(&u).Error
}

func (u UserSettingsHistory) Save() error {
	return database.Save(&u).Error
}

func (u User) GenerateAvatar() error {
	if u.HasAvatar {
		return nil
	}

	err := govatar.GenerateFileFromUsername(govatar.MALE, u.Username, "data/images/"+u.Uuid+"_av.jpeg")
	if err != nil {
		return err
	}

	u.HasAvatar = true
	return u.Save()
}

/*
	View Models
*/

type ViewUser struct {
	*User
	Balance             float64
	InviterUsername     string
	RegistrationDateStr string
	LastLoginDateStr    string
	IsOnline            bool
}

func (u User) ViewUser(lang string) ViewUser {

	user := ViewUser{
		User:                &u,
		RegistrationDateStr: util.HumanizeTime(u.RegistrationDate, lang),
	}

	if u.LastLoginDate != nil {
		user.LastLoginDateStr = util.HumanizeTime(*u.LastLoginDate, lang)
		user.IsOnline = u.LastLoginDate.After(time.Now().Add(-onlineDuration))
	}

	return user
}

func (users Users) ViewUsers(lang string) []ViewUser {
	viewUsers := []ViewUser{}
	for _, user := range users {
		viewUsers = append(viewUsers, user.ViewUser(lang))
	}
	return viewUsers
}

type ViewExtendedUser struct {
	*ExtendedUser
	Balance             float64
	InviterUsername     string
	RegistrationDateStr string
	LastLoginDateStr    string
	IsOnline            bool
}

func (u ExtendedUser) ViewExtendedUser(lang string) ViewExtendedUser {

	user := ViewExtendedUser{
		ExtendedUser:        &u,
		RegistrationDateStr: util.HumanizeTime(u.RegistrationDate, lang),
	}

	if u.LastLoginDate != nil {
		user.LastLoginDateStr = util.HumanizeTime(*u.LastLoginDate, lang)
		user.IsOnline = u.LastLoginDate.After(time.Now().Add(-onlineDuration))
	}

	return user
}

func (users ExtendedUsers) ViewExtendedUsers(lang string) []ViewExtendedUser {
	viewUsers := []ViewExtendedUser{}
	for _, user := range users {
		viewUsers = append(viewUsers, user.ViewExtendedUser(lang))
	}
	return viewUsers
}

type ViewSeller struct {
	*Seller
	ViewItems           []ViewItem
	ItemCategories      []ItemCategory
	ViewRatingReviews   []ViewRatingReview
	InviterUsername     string
	RegistrationDateStr string
	LastLoginDateStr    string
	NumberOfItems       int
	Score               string
	ScoreFloat          float64
	IsOnline            bool
	MultisigEnabled     bool
	LongDescriptionHTML template.HTML
	NumberOfSales       int
	NumberOfSalesStr    string
	SalesVolume         float64
	SalesVolumeStr      string
	NumberOfPurchases   int
	PurchaseVolume      float64
	// Bitcoin Wallets
	BitcoinBalance       apis.BTCWalletBalance
	BitcoinWallets       UserBitcoinWallets
	BitcoinWallet        UserBitcoinWallet
	BitcoinWalletActions []UserBitcoinWalletAction
	// Ethereum Wallets
	EthereumBalance       apis.ETHWalletBalance
	EthereumWallets       UserEthereumWallets
	EthereumWallet        UserEthereumWallet
	EthereumWalletActions []UserEthereumWalletAction
	// Bitcoin Cash Wallets
	BitcoinCashBalance       apis.BCHWalletBalance
	BitcoinCashWallets       UserBitcoinCashWallets
	BitcoinCashWallet        UserBitcoinCashWallet
	BitcoinCashWalletActions []UserBitcoinCashWalletAction
}

func (s Seller) ViewSeller(lang string) ViewSeller {

	var lastLoginDate string = "?"
	if s.LastLoginDate != nil {
		if lang == "ru" {
			lastLoginDate = util.HumanizeTimeRU(*s.LastLoginDate)
		} else {
			lastLoginDate = humanize.Time(*s.LastLoginDate)
		}
	}

	score := s.Score()

	viewSeller := ViewSeller{
		LastLoginDateStr:    lastLoginDate,
		LongDescriptionHTML: template.HTML(userHtmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(s.User.LongDescription))))),
		MultisigEnabled:     (s.BitcoinMultisigPublicKey != "" && s.Premium),
		RegistrationDateStr: humanize.Time(s.RegistrationDate),
		Score:               humanize.Ftoa(score),
		ScoreFloat:          score,
		Seller:              &s,
	}

	if s.LastLoginDate != nil {
		viewSeller.IsOnline = s.LastLoginDate.After(time.Now().Add(-onlineDuration))
	}

	for _, item := range s.Items {
		nc := item.ItemCategory
		itemExists := false
		for _, ec := range viewSeller.ItemCategories {
			if ec.ID == nc.ID {
				itemExists = true
				break
			}
		}
		if !itemExists {
			viewSeller.ItemCategories = append(viewSeller.ItemCategories, nc)
		}
	}

	for _, review := range viewSeller.RatingReviews {
		viewSeller.ViewRatingReviews = append(viewSeller.ViewRatingReviews, review.ViewRatingReview(lang))
	}

	if lang == "ru" {
		viewSeller.RegistrationDateStr = util.HumanizeTimeRU(s.RegistrationDate)
	}

	inviter := s.Iniviter()
	if inviter != nil {
		viewSeller.InviterUsername = inviter.Username
	}

	return viewSeller
}

type ViewSellers []ViewSeller

// Create views and other representatives
func setupUserViews() {
	database.Exec("DROP VIEW IF EXISTS v_users CASCADE;")
	database.Exec(`
		CREATE VIEW v_users AS (
             select
				u1.*,
				u2.username as inviter_username,
				(select count(*) from users u3 where u3.inviter_uuid=u1.uuid) as inviter_count,
				COALESCE(vubb.balance, 0) as bitcoin_balance,
				COALESCE(vubb.unconfirmed_balance, 0) as bitcoin_unconfirmed_balance,
				COALESCE(vueb.balance, 0) as ethereum_balance,
				COALESCE(vubcb.balance, 0) as bitcoin_cash_balance,
				COALESCE(vubcb.unconfirmed_balance, 0) as bitcoin_cash_unconfirmed_balance,
				COALESCE(vst.number_of_messages, 0) as number_of_support_messages,
				COALESCE(vst.last_message_by_staff, false) as last_message_by_staff,
				vst.last_updated_at as last_message_datetime,
				vst.support_user_username
     		from
				users u1
     		left outer join
				users u2 on u2.uuid=u1.inviter_uuid
     		left outer join
				v_user_bitcoin_wallet_balances vubb on vubb.user_uuid=u1.uuid
     		left outer join
				v_user_ethereum_wallet_balances vueb on vueb.user_uuid=u1.uuid
     		left outer join
				v_user_bitcoin_cash_wallet_balances vubcb on vubcb.user_uuid=u1.uuid
			left outer join
				v_support_threads vst on vst.sender_uuid = u1.uuid
    );`)
}

func refreshUsesMaterializedView() {
	database.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY vm_users ;")
}
