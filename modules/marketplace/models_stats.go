package marketplace

import (
	"fmt"
	"time"
)

type StatsItem struct {
	Date                           time.Time
	DateStr                        string
	Year                           int
	WeekNumber                     int
	NumberOfUsers                  int
	NumberOfUsersDelta             int
	NumberOfVendors                int
	NumberOfVendorsDelta           int
	NumberOfItems                  int
	NumberOfItemsDelta             int
	NumberOfSupportMessages        int
	NumberOfSupportMessagesDelta   int
	NumberOfDisputes               int
	NumberOfDisputesDelta          int
	NumberOfBTCTransactionsCreated int
	NumberOfETHTransactionsCreated int
	NumberOfBCHTransactionsCreated int
	BTCTradeAmount                 float64
	ETHTradeAmount                 float64
	BCHTradeAmount                 float64
}

func GetMarketplaceStats(dt time.Time) []StatsItem {
	database.Exec("REFRESH MATERIALIZED VIEW vm_current_bitcoin_transaction_statuses")
	database.Exec("REFRESH MATERIALIZED VIEW vm_current_bitcoin_cash_transaction_statuses")
	database.Exec("REFRESH MATERIALIZED VIEW vm_current_ethereum_transaction_statuses")

	d, _ := time.ParseDuration("168h")

	btcTable := "vm_current_bitcoin_transaction_statuses"
	ethTable := "vm_current_ethereum_transaction_statuses"
	bchTable := "vm_current_bitcoin_cash_transaction_statuses"

	statsItems := []StatsItem{}
	for {
		if dt.After(time.Now()) {
			break
		}
		year, week := dt.ISOWeek()
		from := dt
		to := dt.Add(d)
		statItem := StatsItem{
			NumberOfUsers:                  CountUsers(&to),
			NumberOfVendors:                CountVendors(&to),
			NumberOfBTCTransactionsCreated: CountTransactionsInitiatedInPeriod(btcTable, from, to),
			BTCTradeAmount:                 SumAmountTransactionsInitiatedInPeriod(btcTable, from, to),
			NumberOfETHTransactionsCreated: CountTransactionsInitiatedInPeriod(ethTable, from, to),
			ETHTradeAmount:                 SumAmountTransactionsInitiatedInPeriod(ethTable, from, to),
			NumberOfBCHTransactionsCreated: CountTransactionsInitiatedInPeriod(bchTable, from, to),
			BCHTradeAmount:                 SumAmountTransactionsInitiatedInPeriod(bchTable, from, to),
			Date:                           from,
			Year:                           year,
			WeekNumber:                     week,
			DateStr:                        dt.Format("2006-01-02"),
		}
		statsItems = append(statsItems, statItem)
		dt = dt.Add(d)
	}
	return statsItems
}

// Stats tables are different for various cryptocurrencies
// vm_current_bitcoin_transaction_statuses
// vm_current_ethereum_transaction_statuses

func CountTransactionsInitiatedInPeriod(table string, from, to time.Time) int {
	var count int

	database.
		Table(table).
		Where("amount > 0").
		Where("current_status <> ? AND current_status <> ? AND current_status <> ?", "FAILED", "PENDING", "CANCELLED").
		Where("created_at >= ? and created_at < ?", from, to).
		Count(&count)

	return count
}

func CountTransactionsLastUpdatedWithStatusInPeriod(table string, from, to time.Time, status string) int {
	var count int

	database.
		Table(table).
		Where("amount > 0").
		Where("current_status = ?", status).
		Where("updated_at >= ? and updated_at < ?", from, to).
		Count(&count)

	return count
}

func SumAmountTransactionsInitiatedInPeriod(table string, from, to time.Time) float64 {

	var sum struct{ Sum float64 }

	database.
		Table(table).
		Select("COALESCE(SUM(amount), 0) as sum").
		Where("amount > 0").
		Where("current_status <> ? AND current_status <> ? AND current_status <> ?", "FAILED", "PENDING", "CANCELLED").
		Where("created_at >= ? and created_at < ?", from, to).
		Scan(&sum)

	return sum.Sum
}

type StaffSupportTicketsStatsItem struct {
	ResolverUsername string
	CurrentStatus    string
	TicketCount      int
}

func StaffSupportTicketsResolutionStats(interval string) ([]StaffSupportTicketsStatsItem, error) {
	var (
		query = fmt.Sprintf(`
select 
	u.username as resolver_username,
	current_status,
	count(*) as ticket_count
from 
	v_current_support_ticket_statuses ts 
join 
	users u on ts.resolver_user_uuid=u.uuid
where 
	u.is_staff=true AND ts.updated_at BETWEEN now() - interval '%s' AND now() 
group by 
	u.username, current_status;
`, interval)
		results = []StaffSupportTicketsStatsItem{}
		err     = database.Raw(query).Find(&results).Error
	)
	return results, err
}

type CountStatsItem struct {
	Count int
}

func NumberOfNewUsersStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select
	count(*) as count 
from
	users 
where 
	registration_date between now() - interval '%s' and now();
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfActiveUsersStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select
	count(*) as count 
from
	users 
where 
	last_login_date between now() - interval '%s' and now();
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfNewTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfCompletedTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()  AND max_status = 'COMPLETED'
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfReleasedTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()  AND max_status = 'RELEASED'
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfFailedTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()  AND max_status = 'FAILED'
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfCancelledTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()  AND max_status = 'CANCELLED'
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

func NumberOfFrozenTransactionsStats(interval string) int {
	var (
		query = fmt.Sprintf(`
select 
	count(*) as count
from 
	v_transaction_statuses where min_timestamp BETWEEN now() - interval '%s' and now()  AND max_status = 'FROZEN'
`, interval)
		count = CountStatsItem{}
	)

	database.Raw(query).Scan(&count)
	return count.Count
}

/*
	Cache
*/

func CacheGetMarketplaceStats(dt time.Time) []StatsItem {
	key := "stats-" + dt.String()
	cStats, _ := gc.Get(key)
	if cStats == nil {
		stats := GetMarketplaceStats(dt)
		gc.Set(key, stats)
		return stats
	}
	return cStats.([]StatsItem)
}
