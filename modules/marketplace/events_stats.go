package marketplace

import (
	"fmt"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventStaffSupportTicketsResolutionStats(interval string, sItems []StaffSupportTicketsStatsItem) {
	var (
		text = fmt.Sprintf(`
# %s Support Ticket Resolution Statistics

Support Agent | Ticket Status | Number Of Tickets
--- | --- | ---
`, interval)
	)
	for _, si := range sItems {
		text += fmt.Sprintf("%s | %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}

	hook := ""
	switch interval {
	case "1 day":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsDaily
	case "7 days":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsWeekly
	case "1 hour":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsHourly
	}

	apis.PostMattermostEvent(hook, text)
}

func EventUsersStats(interval string, newUsers, activeUsers int) {
	var (
		text = fmt.Sprintf(`
# %s Users Stats

New Users | Active Users
--- | ---
%d | %d

`, interval, newUsers, activeUsers)
	)

	hook := ""
	switch interval {
	case "1 day":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsDaily
	case "7 days":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsWeekly
	case "1 hour":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsHourly
	}

	apis.PostMattermostEvent(hook, text)
}

func EventTransactionsStats(
	interval string,
	numberOfNewTransactions,
	numberOfCompletedTransactions,
	numberOfReleasedTransactions,
	numberOfFailedTransactions,
	numberOfCancelledTransactions,
	numberOfFrozenTransactions int,
) {
	var (
		text = fmt.Sprintf(`
# %s Transactions Stats

New | Completed | Released | Failed | Cancelled | Frozen
--- | --- | --- | --- | --- | --- 
%d | %d | %d | %d | %d | %d
`,
			interval,
			numberOfNewTransactions,
			numberOfCompletedTransactions,
			numberOfReleasedTransactions,
			numberOfFailedTransactions,
			numberOfCancelledTransactions,
			numberOfFrozenTransactions,
		)
	)

	hook := ""
	switch interval {
	case "1 day":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsDaily
	case "7 days":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsWeekly
	case "1 hour":
		hook = MARKETPLACE_SETTINGS.MattermostIncomingHookStatsHourly
	}

	apis.PostMattermostEvent(hook, text)
}
