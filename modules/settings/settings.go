package settings

import (
	"encoding/json"
	"io/ioutil"
)

type MarketplaceSettings struct {
	CompletedDuration                    string  `json:"completed_duration"`
	CookieEncryptionSalt                 string  `json:"cookie_encryption_salt"`
	CooloffPeriod                        int     `json:"cooloff_period"`
	FeedSize                             int     `json:"feed_size"`
	FreeCommission                       float64 `json:"free_commission"`
	GRAMSApiKey                          string  `json:"grams_api_key"`
	ItemsPerPage                         int     `json:"items_per_page"`
	MattermostIncomingHookSupport        string  `json:"mattermost_incoming_hook_support"`
	MattermostIncomingHookAuthentication string  `json:"mattermost_incoming_hook_authentication"`
	MattermostIncomingHookMessageboard   string  `json:"mattermost_incoming_hook_messageboard"`
	MattermostIncomingHookTrustedVendors string  `json:"mattermost_incoming_hook_trusted_vendors"`
	MattermostIncomingHookTransactions   string  `json:"mattermost_incoming_hook_transactions"`
	MattermostIncomingHookShoutbox       string  `json:"mattermost_incoming_hook_shoutbox"`
	MattermostIncomingHookItems          string  `json:"mattermost_incoming_hook_items"`
	MattermostIncomingHookDisputes       string  `json:"mattermost_incoming_hook_disputes"`
	MattermostIncomingHookStatsDaily     string  `json:"mattermost_incoming_hook_stats_daily"`
	MattermostIncomingHookStatsHourly    string  `json:"mattermost_incoming_hook_stats_hourly"`
	MattermostIncomingHookStatsWeekly    string  `json:"mattermost_incoming_hook_stats_weekly"`
	OnlineDuration                       string  `json:"online_duration"`
	PendingDuration                      string  `json:"pending_duration"`
	PostgresConnectionString             string  `json:"postgres_connection_string"`
	PremiumCommission                    float64 `json:"premium_commission"`
	PremiumPlusCommission                float64 `json:"premium_plus_commission"`
	AdvertisingCost                      float64 `json:"advertising_cost"`
	Sitename                             string  `json:"sitename"`
	Host                                 string  `json:"host"`
	Port                                 string  `json:"port"`
	PaymentGate                          string  `json:"paymentgate"`
	Debug                                bool    `json:"debug"`
	SingleMode                           bool    `json:"single_mode"`
}

var (
	settings MarketplaceSettings
)

func GetSettings() MarketplaceSettings {

	settings = loadmarketplacesettings()

	if settings.Host == "" {
		settings.Host = "127.0.0.1"
	}

	if settings.Port == "" {
		settings.Port = "8081"
	}

	if settings.PaymentGate == "" {
		settings.PaymentGate = "http://127.0.0.1:8083"
	}

	return settings
}

func loadmarketplacesettings() MarketplaceSettings {

	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		panic(err)
	}

	return settings
}
