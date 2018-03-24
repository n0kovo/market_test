package marketplace

import (
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

type Context struct {
	*util.Context
	// Localization
	Localization Localization `json:"-"`
	// General
	CanEdit   bool   `json:"can_edit,omitempty"`
	CaptchaId string `json:"captcha_id,omitempty"`
	Error     string `json:"error,omitempty"`
	//Mode marketplace
	IsSingleMode bool `json:"-"`
	// Misc
	Pgp                 string                `json:"pgp,omitempty"`
	UserSettingsHistory []UserSettingsHistory `json:"user_settings_history,omitempty"`
	Language            string                `json:"language,omitempty"`
	// Paging & sorting
	SelectedPage  int    `json:"selected_page,omitempty"`
	Pages         []int  `json:"-,omitempty"`
	Page          int    `json:"page,omitempty"`
	NumberOfPages int    `json:"number_of_pages,omitempty"`
	Query         string `json:"query,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	// Static Pages
	StaticPage  StaticPage   `json:"-,omitempty"`
	StaticPages []StaticPage `json:"-,omitempty"`
	// Messageboard
	IsPrivateMessage bool `json:"is_private_message,omitempty"`
	// ReadOnlyThread   bool
	ItemCategories []ItemCategory `json:"-,omitempty"`
	ItemCategory   ItemCategory   `json:"-,omitempty"`
	// Menu
	Categories          []Category     `json:"-,omitempty"`
	Cities              map[string]int `json:"-,omitempty"`
	City                string         `json:"city,omitempty"`
	GeoCities           []City         `json:"geo_cities,omitempty"`
	CityID              int            `json:"city_id,omitempty"`
	Countries           []Country      `json:"countries,omitempty"`
	Quantity            int            `json:"quantity,omitempty"`
	SelectedPackageType string         `json:"selected_package_type,omitempty"`
	SelectedSection     string         `json:"-,omitempty"`
	SelectedSectionID   int            `json:"-,omitempty"`
	SelectedStatus      string         `json:"selected_status,omitempty"`
	ShippingFrom        string         `json:"shipping_from,omitempty"`
	ShippingFromList    []string       `json:"shipping_from_list,omitempty"`
	ShippingTo          string         `json:"shipping_to,omitempty"`
	ShippingToList      []string       `json:"shipping_to_list,omitempty"`
	Account             string         `json:"account,omitempty"`
	// Categories
	Category       string `json:"category,omitempty"`
	SubCategory    string `json:"sub_category,omitempty"`
	SubSubCategory string `json:"sub_sub_category,omitempty"`
	CategoryID     int    `json:"category_id,omitempty"`
	// Items page
	GroupPackages      []GroupPackage            `json:"group_packages,omitempty"`
	GroupPackagesTable map[GroupPackage]Packages `json:"group_packages_table,omitempty"`
	GroupAvailability  GroupPackage              `json:"-,omitempty"`
	NumberOfItems      int                       `json:"number_of_items,omitempty"`
	PackageCurrency    string                    `json:"package_currency,omitempty"`
	PackagePrice       string                    `json:"package_price,omitempty"`
	// Transactions page
	PendingCount   int `json:"pending_count,omitempty"`
	FailedCount    int `json:"failed_count,omitempty"`
	ReleasedCount  int `json:"released_count,omitempty"`
	CompletedCount int `json:"completed_count,omitempty"`
	AllCount       int `json:"all_count,omitempty"`
	// Models
	ExtendedUsers        []ExtendedUser        `json:"-,omitempty"`
	Invitation           Invitation            `json:"-,omitempty"`
	Invitations          []Invitation          `json:"-,omitempty"`
	Item                 Item                  `json:"-,omitempty"`
	Items                Items                 `json:"-,omitempty"`
	Package              Package               `json:"-,omitempty"`
	Packages             Packages              `json:"-,omitempty"`
	Seller               Seller                `json:"-,omitempty"`
	Sellers              Sellers               `json:"-,omitempty"`
	Thread               Thread                `json:"-,omitempty"`
	Threads              []Thread              `json:"-,omitempty"`
	Transaction          Transaction           `json:"-,omitempty"`
	Transactions         []Transaction         `json:"-,omitempty"`
	MessageboardSections []MessageboardSection `json:"-,omitempty"`
	MessageboardSection  MessageboardSection   `json:"-,omitempty"`
	RatingReview         RatingReview          `json:"-,omitempty"`
	// View Models
	ViewCurrentTransactionStatuses []ViewCurrentTransactionStatus `json:"-"`
	ViewExtendedUsers              []ViewExtendedUser             `json:"-"`
	ViewFeedItems                  []ViewFeedItem                 `json:"-"`
	ViewInvitation                 ViewInvitation                 `json:"-"`
	ViewItem                       ViewItem                       `json:"-"`
	ViewItems                      []ViewItem                     `json:"-"`
	ViewMessage                    ViewMessage                    `json:"-"`
	ViewMessages                   []ViewMessage                  `json:"-"`
	ViewPackage                    ViewPackage                    `json:"-"`
	ViewPackages                   []ViewPackage                  `json:"-"`
	ViewSeller                     ViewSeller                     `json:"-"`
	ViewSellers                    []ViewSeller                   `json:"-"`
	ViewThread                     ViewThread                     `json:"-"`
	ViewThreads                    []ViewThread                   `json:"-"`
	ViewTransaction                ViewTransaction                `json:"-"`
	ViewTransactions               []ViewTransaction              `json:"-"`
	ViewUser                       ViewUser                       `json:"-"`
	ViewUsers                      []ViewUser                     `json:"-"`
	// Stats
	NumberOfDailyTransactions     int `json:"-"`
	NumberOfMonthlyTransactions   int `json:"-"`
	NumberOfPrivateMessages       int `json:"-"`
	NumberOfSupportMessages       int `json:"-"`
	NumberOfTransactions          int `json:"-"`
	NumberOfUnreadPrivateMessages int `json:"-"`
	NumberOfUnreadSupportMessages int `json:"-"`
	NumberOfWeeklyTransactions    int `json:"-"`
	NumberOfDisputes              int `json:"-"`
	// Admin Stats
	NumberOfUsers              int         `json:"-"`
	NumberOfVendors            int         `json:"-"`
	NumberOfVendorsFree        int         `json:"-"`
	NumberOfVendorsPremium     int         `json:"-"`
	NumberOfVendorsPremiumPlus int         `json:"-"`
	NumberOfNewUsers           int         `json:"-"`
	NumberOfActiveUsers        int         `json:"-"`
	NumberOfWeeklyActiveUsers  int         `json:"-"`
	NumberOfOnlineUsers        int         `json:"-"`
	NumberOfMonthlyActiveUsers int         `json:"-"`
	NumberOfInvitedUsers       int         `json:"-"`
	StatsItems                 []StatsItem `json:"-"`
	// Auth
	SecretText string `json:"secret_text,omitempty"`
	InviteCode string `json:"invite_code,omitempty"`
	// Bitcoin Wallets
	UserBitcoinBalance       apis.BTCWalletBalance     `json:"-"`
	UserBitcoinWallets       UserBitcoinWallets        `json:"-"`
	UserBitcoinWallet        UserBitcoinWallet         `json:"-"`
	UserBitcoinWalletActions []UserBitcoinWalletAction `json:"-"`
	// Ethereum Wallets
	UserEthereumBalance       apis.ETHWalletBalance      `json:"-"`
	UserEthereumWallets       UserEthereumWallets        `json:"-"`
	UserEthereumWallet        UserEthereumWallet         `json:"-"`
	UserEthereumWalletActions []UserEthereumWalletAction `json:"-"`
	// Bitcoin Cash Wallets
	UserBitcoinCashBalance       apis.BCHWalletBalance         `json:"-"`
	UserBitcoinCashWallets       UserBitcoinCashWallets        `json:"-"`
	UserBitcoinCashWallet        UserBitcoinCashWallet         `json:"-"`
	UserBitcoinCashWalletActions []UserBitcoinCashWalletAction `json:"-"`
	// Referrals
	ReferralPayments []ReferralPayment `json:"-"`
	//Dispute
	Dispute      Dispute      `json:"-"`
	Disputes     []Dispute    `json:"-"`
	DisputeClaim DisputeClaim `json:"-"`
	// Support
	SupportThreads          []SupportThread          `json:"-"`
	ViewMessageboardThreads []ViewMessageboardThread `json:"-"`
	ViewSupportTicket       ViewSupportTicket        `json:"-"`
	ViewSupportTickets      []ViewSupportTicket      `json:"-"`
	// New Items List page
	ViewAvailableItems []ViewAvailableItem `json:"available_items,omitempty"`
	ViewVendors        []ViewVendor        `json:"-"`
	// Currency Rates
	CurrencyRates map[string]map[string]float64 `json:"-"`
	USDBTCRate    float64                       `json:"-"`
	// Wallet page
	BTCFee            float64                 `json:"-"`
	BCHFee            float64                 `json:"-"`
	Amount            float64                 `json:"-"`
	Address           string                  `json:"-"`
	Description       string                  `json:"-"`
	BTCPaymentResult  apis.BTCPaymentResult   `json:"-"`
	BCHPaymentResult  apis.BCHPaymentResult   `json:"-"`
	ETHPaymentResults []apis.ETHPaymentResult `json:"-"`
	// Advertising
	Advertisings    []Advertising `json:"-"`
	AdvertisingCost float64       `json:"-"`
	// ApiSession
	APISession *APISession `json:"api_session,omitempty"`
}
