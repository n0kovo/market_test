package marketplace

import (
	"fmt"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type MessageboardThread struct {
	Uuid                  string
	LastMessageUsername   string
	MessageboardSectionId int
	NumberOfMessages      int
	Pages                 []int
	SenderIsAdmin         bool
	SenderIsPremium       bool
	SenderIsPremiumPlus   bool
	SenderIsStaff         bool
	SenderUsername        string
	Title                 string

	CreatedAt     time.Time
	LastUpdatedAt time.Time
	LastRead      time.Time
}

type ViewMessageboardThread struct {
	*MessageboardThread
	IsRead           bool
	CreatedAtStr     string
	LastUpdatedAtStr string
}

func (mt MessageboardThread) ViewMessageboardThread(lang string) ViewMessageboardThread {
	vmbt := ViewMessageboardThread{
		MessageboardThread: &mt,
		CreatedAtStr:       util.HumanizeTime(mt.CreatedAt, lang),
		LastUpdatedAtStr:   util.HumanizeTime(mt.LastUpdatedAt, lang),
	}
	return vmbt
}

type MessageboardThreads []MessageboardThread

func (msbts MessageboardThreads) ViewMessageboardThreads(lang string) []ViewMessageboardThread {
	var vmsbts []ViewMessageboardThread
	for _, msbt := range msbts {
		vmsbt := msbt.ViewMessageboardThread(lang)
		vmsbts = append(vmsbts, vmsbt)
	}

	return vmsbts
}

func FindMessageboardThreads(sectionId int) []MessageboardThread {
	threads := []MessageboardThread{}

	query := database.
		Table("v_messageboard_threads").
		Where("messageboard_section_id=?", sectionId)
	query.Find(&threads)

	return threads
}

func FindMessageboardThreadsForUserUuid(sectionId, page, pageSize int, userUuid string) []MessageboardThread {
	threads := []MessageboardThread{}

	database.
		Table("v_messageboard_threads").
		Select("v_messageboard_threads.*, thread_perusal_statuses.last_read_date as last_read").
		Joins(fmt.Sprintf("left outer join thread_perusal_statuses on v_messageboard_threads.uuid=thread_perusal_statuses.thread_uuid AND thread_perusal_statuses.user_uuid = '%s'", userUuid)).
		Where("messageboard_section_id=?", sectionId).
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&threads)

	return threads
}

func FindStaffMessageboardThreadsForUserUuid(sectionId, page, pageSize int, userUuid string) []MessageboardThread {
	threads := []MessageboardThread{}

	database.
		Table("v_staff_messageboard_threads").
		Select("v_staff_messageboard_threads.*, thread_perusal_statuses.last_read_date as last_read").
		Joins(fmt.Sprintf("left outer join thread_perusal_statuses on v_staff_messageboard_threads.uuid=thread_perusal_statuses.thread_uuid AND thread_perusal_statuses.user_uuid = '%s'", userUuid)).
		// Where("messageboard_section_id=?", sectionId).
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&threads)

	return threads
}

func CountMessageboardThreads(sectionId int) int {
	count := 0

	query := database.Table("v_messageboard_threads").Where("messageboard_section_id=?", sectionId)
	query.Count(&count)

	return count
}

func CountStaffMessageboardThreads(sectionId int) int {
	count := 0

	query := database.Table("v_staff_messageboard_threads")
	query.Count(&count)

	return count
}

func setupMessageboardThreadsViews() {
	database.Exec("DROP VIEW IF EXISTS v_messageboard_threads CASCADE;")
	database.Exec(`
		CREATE VIEW v_messageboard_threads AS (
			SELECT 
				v_threads.uuid,
				v_threads.title,
				v_threads.created_at_timestamp as created_at,
				v_threads.last_updated as last_updated_at,
				v_threads.number_of_messages,
				v_threads.messageboard_section_id,
				u1.username as sender_username,
				u2.username as last_message_username,
				u1.is_admin as sender_is_admin,
				u1.premium as sender_is_premium, 
				u1.premium_plus as sender_is_premium_plus,
				u1.is_staff as sender_is_staff
			FROM v_threads
			JOIN users u1 on u1.uuid=sender_user_uuid
			JOIN messages m on m.uuid=last_message_uuid
			JOIN users u2 on u2.uuid=m.sender_user_uuid
			WHERE v_threads.type='messageboard'
			AND u1.last_login_date >= (now() - interval '21 day')
			ORDER BY u1.is_admin DESC, u1.premium_plus DESC, u1.premium DESC, last_updated_at DESC
	);`)
}

func setupStaffMessageboardThreadsViews() {
	database.Exec("DROP VIEW IF EXISTS v_staff_messageboard_threads CASCADE;")
	database.Exec(`
		CREATE VIEW v_staff_messageboard_threads AS (
			SELECT 
				v_threads.uuid,
				v_threads.title,
				v_threads.created_at_timestamp as created_at,
				v_threads.last_updated as last_updated_at,
				v_threads.number_of_messages,
				v_threads.messageboard_section_id,
				u1.username as sender_username,
				u2.username as last_message_username,
				u1.is_admin as sender_is_admin,
				u1.premium as sender_is_premium, 
				u1.premium_plus as sender_is_premium_plus,
				u1.is_staff as sender_is_staff
			FROM v_threads
			JOIN users u1 on u1.uuid=sender_user_uuid
			JOIN messages m on m.uuid=last_message_uuid
			JOIN users u2 on u2.uuid=m.sender_user_uuid
			WHERE v_threads.type='staff_messageboard'
			AND u1.last_login_date >= (now() - interval '21 day')
			ORDER BY u1.is_admin DESC, u1.premium_plus DESC, u1.premium DESC, last_updated_at DESC
	);`)
}
