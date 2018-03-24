package marketplace

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/settings"
)

var (
	database *gorm.DB
)

func SyncModels() {
	database.AutoMigrate(
		&User{},
		&Item{},
		&PackagePrice{},
		&Package{},
		&BitcoinTransaction{},
		&BitcoinCashTransaction{},
		&EthereumTransaction{},
		&TransactionStatus{},
		&Reservation{},
		&Transaction{},
		&Message{},
		&ThreadPerusalStatus{},
		&SupportTicket{},
		&SupportTicketStatus{},
		&ShippingOption{},
		&ShippingStatus{},
		&UserSettingsHistory{},
		&Invitation{},
		&FeedItem{},
		&ItemCategory{},
		&UserBitcoinWallet{},
		&UserBitcoinWalletBalance{},
		&UserBitcoinWalletAction{},
		&UserBitcoinCashWallet{},
		&UserBitcoinCashWalletBalance{},
		&UserBitcoinCashWalletAction{},
		&UserEthereumWallet{},
		&UserEthereumWalletBalance{},
		&UserEthereumWalletAction{},
		&PaymentReceipt{},
		&MessageboardSection{},
		&RatingReview{},
		&ReferralPayment{},
		&Dispute{},
		&DisputeClaim{},
		&City{},
		&Country{},
		&Advertising{},
		&APISession{},
	)

	// drop all views

	database.Exec(`
	SELECT 
	'DROP VIEW ' || table_name || ';'
	FROM information_schema.views
	WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
	AND table_name !~ '^pg_';
	`)

	createViews := func() {
		// wallets & balances
		setupUserBitcoinBalanceViews()
		setupUserBitcoinCashBalanceViews()
		setupUserEthereumBalanceViews()

		// messageboard & messages
		setupThreadsViews()
		setupSupportThreadsViews()
		setupMessageboardThreadsViews()
		setupStaffMessageboardThreadsViews()

		// transcations
		setupTransactionStatusesView()

		// users
		setupUserViews()
		setupVendorTxStatsViews()
		setupItemTxStatsViews()

		// items & packages, categories
		setupCategoriesViews()
		setupPackagesView()
		setupAvailableItemsView()

		// tickets
		setupSupportTicketViews()
	}

	createViews()
}

func init() {
	var err error

	database, err = gorm.Open("postgres", MARKETPLACE_SETTINGS.PostgresConnectionString)
	if err != nil {
		panic(err)
	}
	database.DB().SetMaxOpenConns(8)
	database.DB().Ping()

	if settings.GetSettings().Debug {
		database.LogMode(true)
	}

}
