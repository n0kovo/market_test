package marketplace

import (
	"time"
)

/*
	Models
*/

type SupportThread struct {
	Uuid                string
	CreatedAt           time.Time
	LastMessageByStaff  bool
	LastMessageUsername string
	LastMessageUuid     string
	LastUpdatedAt       time.Time
	NumberOfMessages    int
	SenderUsername      string
	SenderUuid          string
}

func FindSupportThreads(lastMessageByStaff *bool) []SupportThread {
	threads := []SupportThread{}

	query := database.Table("v_support_threads")
	if lastMessageByStaff != nil {
		query = query.Where("last_message_by_staff=?", lastMessageByStaff)
	}

	query.Find(&threads)

	return threads
}

func GetSupportThread(user User, createIfNotExists bool) (*Thread, error) {
	threadUuid := "support-" + user.Uuid
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "support" {
		return thread, nil
	}
	return CreateThread(
		"support",
		threadUuid,
		"Support thread @"+user.Username,
		"",
		&user,
		nil,
		createIfNotExists,
	)
}

func setupSupportThreadsViews() {
	database.Exec("DROP VIEW IF EXISTS v_support_threads CASCADE;")
	database.Exec(`
		CREATE VIEW v_support_threads AS (
			SELECT 
				v_threads.uuid,
				v_threads.created_at_timestamp as created_at,
				v_threads.last_updated as last_updated_at,
				v_threads.number_of_messages,
				u1.username as sender_username,
				u1.uuid as sender_uuid,
				u2.username as last_message_username,
				u2.username as last_message_uuid,
				(u2.is_admin or u2.is_staff) as last_message_by_staff,
                u3.username as support_user_username
			FROM v_threads
			JOIN users u1 on u1.uuid=sender_user_uuid
			JOIN messages m on m.uuid=last_message_uuid
			JOIN users u2 on u2.uuid=m.sender_user_uuid
            LEFT OUTER JOIN users u3 on u3.uuid=u1.supporter_uuid
			WHERE v_threads.type='support'
			AND u1.last_login_date >= (now() - interval '21 day')
			ORDER BY last_updated_at DESC
	);`)
}
