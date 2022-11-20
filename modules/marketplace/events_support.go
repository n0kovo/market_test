package marketplace

import (
	"fmt"
	"strings"

	"github.com/n0kovo/market_test/modules/apis"
)

func EventNewSupportMessage(message *Message) {
	sender, _ := FindUserByUuid(message.SenderUserUuid, false)
	marketUrl := MARKETPLACE_SETTINGS.Sitename
	thread, _ := FindThreadByUuid(message.ParentUuid)
	text := fmt.Sprintf("[@%s](%s/user/%s) has posted in [support thread @%s](%s/staff/users/support/%s):\n> %s",
		sender,
		marketUrl,
		sender,

		thread.SenderUser.Username,
		marketUrl,
		thread.SenderUser.Username,

		strings.Replace(message.Text, "\n", "\n > ", -1),
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookSupport, text)
}

func EventNewSupportTicketMessage(supportTicket *SupportTicket, message *Message) {
	sender, _ := FindUserByUuid(message.SenderUserUuid, false)
	marketUrl := MARKETPLACE_SETTINGS.Sitename
	text := fmt.Sprintf("[@%s](%s/user/%s) has posted in [support %s](%s/support/%s):\n> %s",
		sender,              //---|----|------|---------------------------|---|----------|-------|
		marketUrl,           //--------|------|---------------------------|---|----------|-------|
		sender,              //---------------|---------------------------|---|----------|-------|
		supportTicket.Title, //-------------------------------------------|---|----------|-------|
		marketUrl,           //-----------------------------------------------|----------|-------|
		supportTicket.Uuid,  //----------------------------------------------------------|-------|

		strings.Replace(message.Text, "\n", "\n > ", -1), //-------------------------|
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookSupport, text)
}

func EventNewSupportTicket(supportTicket *SupportTicket) {
	sender, _ := FindUserByUuid(supportTicket.UserUuid, false)
	marketUrl := MARKETPLACE_SETTINGS.Sitename
	text := fmt.Sprintf("[@%s](%s/user/%s) has created [new ticket %s](%s/support/%s):\n> %s",
		sender,              //---|----|------|---------------------------|---|----------|-------|
		marketUrl,           //--------|------|---------------------------|---|----------|-------|
		sender,              //---------------|---------------------------|---|----------|-------|
		supportTicket.Title, //-------------------------------------------|---|----------|-------|
		marketUrl,           //-----------------------------------------------|----------|-------|
		supportTicket.Uuid,  //----------------------------------------------------------|-------|

		strings.Replace(supportTicket.Description, "\n", "\n > ", -1), //------------|
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookSupport, text)
}
