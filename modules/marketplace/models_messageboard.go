package marketplace

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gocraft/web"
	"github.com/russross/blackfriday"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func init() {
	messageboardHtmlPolicy.AllowElements(
		"p",
		"strong",
		"ul",
		"li",
		"i",
		"br",
	)
}

/*
	Models
*/

type MessageboardSection struct {
	ID       int `json:"id" gorm:"primary_key"`
	Priority int `json:"priority"`
	ParentID int `json:"parent_id"`

	Icon string `json:"icon"`
	Flag string `json:"flag"`

	NameEn string `json:"name_en"`
	NameRu string `json:"name_ru"`
	NameDe string `json:"name_de"`
	NameEs string `json:"name_es"`
	NameFr string `json:"name_fr"`

	DescriptionEn string `json:"description_en"`
	DescriptionRu string `json:"description_ru"`
	DescriptionDe string `json:"description_de"`
	DescriptionEs string `json:"description_es"`
	DescriptionFr string `json:"description_fr"`

	Subsections      []MessageboardSection `json:"subsections"`
	NumberOfMessages int                   `json:"number_of_messages" sql:"-"`
	HeadingSection   bool                  `json:"heading_section"`
}

type Message struct {
	Uuid                  string              `form:"uuid" json:"uuid" gorm:"primary_key" sql:"size:1024"`
	SenderUserUuid        string              `json:"sender_uuid" sql:"index""`
	RecieverUserUuid      string              `json:"reciever_uuid" sql:"index"`
	IsReadByReciever      bool                `json:"is_read_by_reciever"`
	Section               string              `form:"section" json:"section" sql:"index"`
	MessageboardSectionID int                 `form:"section_id" json:"section_id" sql:"index"`
	ParentUuid            string              `form:"parent_uuid" json:"parent_uuid" sql:"index"`
	Title                 string              `form:"title" json:"title" sql:"size:140"`
	Text                  string              `form:"text" json:"text" sql:"size:9086"`
	Type                  string              `json:"type" sql:"index"`
	IsEncrypted           bool                `json:"is_encrypted"`
	HasImage              bool                `json:"has_image" sql:"index"`
	CreatedAtTimestamp    time.Time           `json:"created_at" sql:"index"`
	UpdatedAt             *time.Time          `json:"updated_at" sql:"index"`
	DeletedAt             *time.Time          `json:"deleted_at" sql:"index"`
	SenderUser            User                `json:"-"`
	RecieverUser          User                `json:"-"`
	MessageboardSection   MessageboardSection `json:"-"`
}

type Messages []Message

type Thread struct {
	Message
	LastUpdated         time.Time
	LastRead            *time.Time
	ThreadSupportStatus bool
}

type Threads []Thread

type ThreadPerusalStatus struct {
	ThreadUuid   string    `json:"thread_uuid" gorm:"primary_key" sql:"size:1024"`
	UserUuid     string    `json:"user_uuid" gorm:"primary_key" sql:"size:1024"`
	LastReadDate time.Time `json:"last_reade_date" sql:"index"`
}

type ThreadSupportStatuses []ThreadSupportStatus

type ThreadSupportStatus struct {
	ThreadSupportOptionUuid string    `json:"uuid" gorm:"primary_key" sql:"size:1024"`
	UserUuidMark            string    `json:"user_uuid_mark" gorm:"primary_key" sql:"size:1024"`
	DateMark                time.Time `json:"date_mark" sql:"index"`
	MessageUuidMark         string    `json:"message_uuid_mark" sql:"size:1024"`
	IsFixProblem            bool      `json:"is_fix_problem" sql:"index"`
}

/*
	Model Methods
*/

func (message Message) AddImage(r *web.Request) error {

	_, handler, err := r.FormFile("image")

	if handler.Size == 0 {
		return nil
	}

	switch err {
	case nil:
		err = util.SaveImage(r, "image", 2048, message.Uuid)
		if err != nil {
			return err
		}
		message.HasImage = true
		return message.Save()
	case http.ErrMissingFile:
		return nil
	default:
		return err
	}

}

func (m Message) Validate() error {
	if m.Uuid == "" {
		return errors.New("Empty UUID")
	}
	if m.SenderUserUuid == "" && m.Type != "shoutbox" && m.Type != "news" && m.Type != "dispute" && m.Type != "support" {
		return errors.New("Empty User Uuid")
	}
	if m.Text == "" && m.Type == "messageboard" {
		return errors.New("Empty text")
	}
	return nil
}

func (m Thread) Validate() error {
	err := m.Message.Validate()
	if err != nil {
		return err
	}
	if m.Title == "" {
		return errors.New("Empty title")
	}
	if m.Type == "private" && m.RecieverUserUuid == "" {
		return errors.New("Private thread should have reciever")
	}
	return nil
}

func (ts ThreadPerusalStatus) Validate() error {
	if ts.ThreadUuid == "" {
		return errors.New("Empty thread uuid")
	}
	if th, _ := FindThreadByUuid(ts.ThreadUuid); th == nil {
		return errors.New("No such thread")
	}

	if ts.UserUuid == "" {
		return errors.New("Empty thread uuid")
	}
	if u, _ := FindUserByUuid(ts.UserUuid, false); u == nil {
		return errors.New("No such user")
	}

	return nil
}

/*
	Database methods
*/

func (i Message) Remove() error {
	return database.Delete(&i).Error
}

func (itm Message) Save() error {
	err := itm.Validate()
	if err != nil {
		return err
	}
	return itm.SaveToDatabase()
}

func (ts ThreadPerusalStatus) Save() error {
	err := ts.Validate()
	if err != nil {
		return err
	}
	return ts.SaveToDatabase()
}

func (itm Message) SaveToDatabase() error {
	if existing, _ := FindMessageByUuid(itm.Uuid); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

func (itm ThreadPerusalStatus) SaveToDatabase() error {
	if existing, _ := FindThreadPerusalStatus(itm.ThreadUuid, itm.UserUuid); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}
func (tss ThreadSupportStatus) SaveToDatabase() error {

	return database.Save(&tss).Error
}

func (cat MessageboardSection) Validate() error {
	return nil
}

func (cat MessageboardSection) Remove() error {
	return database.Delete(&cat).Error
}

func (cat *MessageboardSection) Save() error {
	err := cat.Validate()
	if err != nil {
		return err
	}
	return cat.SaveToDatabase()
}

func (cat *MessageboardSection) SaveToDatabase() error {
	if existing, _ := FindMessageboardSectionByID(cat.ID); existing == nil {
		cat.ID = rand.Intn(100000)
		return database.Create(cat).Error
	}
	return database.Save(cat).Error
}

/*
	Relations
*/

func (t Thread) Messages(order string) Messages {
	var (
		messages []Message
	)
	q := database.
		Where("messages.parent_uuid=? or messages.uuid=?", t.Uuid, t.Uuid).
		Preload("SenderUser").
		Preload("RecieverUser")
	if order == "ASC" {
		q = q.Order("created_at_timestamp ASC")
	} else {
		q = q.Order("created_at_timestamp DESC")
	}

	q.Find(&messages)

	return messages
}

/*
	Queries
*/

func GetAllMessages() []Message {
	var messages []Message
	database.Find(&messages)
	return messages
}

func GetAllMessageboardSections() []MessageboardSection {
	var messages []MessageboardSection
	database.Find(&messages)
	return messages
}

func FindParentMessageboardSections() []MessageboardSection {
	var sections []MessageboardSection

	database.
		Model(MessageboardSection{}).
		Where("parent_id=0 or parent_id is NULL").
		Order("priority DESC").
		Find(&sections)

	messageboardSectionCounts := CountMessageboardThreadsBySection()

	for i, _ := range sections {
		subsections := FindMessageboardsectionsByParentID(sections[i].ID)
		sections[i].Subsections = subsections
		sections[i].NumberOfMessages = messageboardSectionCounts.CountByID(sections[i].ID)

		for j, _ := range sections[i].Subsections {
			sections[i].Subsections[j].NumberOfMessages = messageboardSectionCounts.CountByID(sections[i].Subsections[j].ID)
		}
	}

	return sections
}

func FindMessageboardsectionsByParentID(parentID int) []MessageboardSection {
	var sections []MessageboardSection

	database.
		Model(MessageboardSection{}).
		Where("parent_id=?", parentID).
		Order("priority DESC").
		Find(&sections)

	return sections
}

func FindThreadsByType(threadType string) Threads {
	var (
		threads Threads
	)

	database.
		Table("v_threads").
		Where("type=?", threadType).
		Order("last_updated ASC").
		Preload("SenderUser").
		Preload("RecieverUser").
		Find(&threads)

	return threads
}

func FindThreadsByActiveUsersByType(threadType string) Threads {
	var (
		threads Threads
	)

	database.
		Table("v_threads").
		Joins("JOIN users on users.uuid=v_threads.sender_user_uuid").
		Where("type=? and users.last_login_date >= (now() - interval '21 day') and users.banned = false", threadType).
		Order("last_updated ASC").
		Preload("SenderUser").
		Preload("RecieverUser").
		Find(&threads)

	return threads
}

type MessageboardSectionThreadCount struct {
	MessageboardSectionID int
	MessageCount          int
}

type MessageboardSectionThreadCounts []MessageboardSectionThreadCount

func (counts MessageboardSectionThreadCounts) CountByID(sectionID int) int {
	count := 0
	for _, s := range counts {
		if s.MessageboardSectionID == sectionID {
			count = s.MessageCount
		}
	}

	return count
}

func CountMessageboardThreadsBySection() MessageboardSectionThreadCounts {
	var results []MessageboardSectionThreadCount

	database.
		Table("v_messageboard_threads").
		Select("v_messageboard_threads.messageboard_section_id as messageboard_section_id, count(*) as message_count").
		Group("v_messageboard_threads.messageboard_section_id").
		Find(&results)

	return MessageboardSectionThreadCounts(results)
}

func GetAllPrivateThreads() Threads {
	threads := FindThreadsByType("private")
	return threads
}

func FindPrivateThreads(user User) Threads {
	var (
		threads Threads
	)

	database.
		Table("v_threads").
		Where(`
			deleted_at IS NULL and
			type=? and 
			(reciever_user_uuid=? OR sender_user_uuid=?)`,
			"private",
			user.Uuid,
			user.Uuid,
		).
		Order("last_updated ASC").
		Preload("SenderUser").
		Preload("RecieverUser").
		Find(&threads)

	return threads
}

func FindPrivateThread(u1, u2 User) Thread {
	var (
		thread Thread
	)

	database.
		Table("v_threads").
		Where(`
			type=? AND (
				(reciever_user_uuid=? AND sender_user_uuid=?) OR
				(reciever_user_uuid=? AND sender_user_uuid=?)
			)`,
			"private",
			u1.Uuid,
			u2.Uuid,
			u2.Uuid,
			u1.Uuid,
		).
		Preload("SenderUser").
		Preload("RecieverUser").
		First(&thread)

	return thread
}

func FindThreadPerusalStatus(threadUuid, userUuid string) (*ThreadPerusalStatus, error) {
	var ts ThreadPerusalStatus

	err := database.
		Where(&ThreadPerusalStatus{
			UserUuid:   userUuid,
			ThreadUuid: threadUuid,
		}).
		First(&ts).Error

	if err != nil {
		return nil, err
	}
	return &ts, err
}

func UpdateThreadPerusalStatus(threadUuid, userUuid string) error {
	ts := ThreadPerusalStatus{
		ThreadUuid:   threadUuid,
		UserUuid:     userUuid,
		LastReadDate: time.Now(),
	}

	return ts.Save()
}
func UpdateThreadSupportStatus(thread Thread) {
	tss := ThreadSupportStatus{
		ThreadSupportOptionUuid: util.GenerateUuid(),
		UserUuidMark:            thread.SenderUserUuid,
		DateMark:                time.Now(),
		MessageUuidMark:         thread.Message.Uuid,
		IsFixProblem:            true,
	}
	tss.SaveToDatabase()
}

func CountPrivateMessages(user User) int {
	// var (
	// 	count = struct{ Sum int }{}
	// )

	// database.
	// 	Select("sum(v_threads.number_of_messages) as sum").
	// 	Table("v_threads").
	// 	Where(`
	// 		deleted_at IS NULL and
	// 		type=? and
	// 		(reciever_user_uuid=? OR sender_user_uuid=?)`,
	// 		"private",
	// 		user.Uuid,
	// 		user.Uuid,
	// 	).
	// 	Scan(&count)

	// return count.Sum

	var (
		count int
	)
	database.
		Table("messages").
		Joins("join users on messages.sender_user_uuid=users.uuid").
		Where(`
messages.type=? and
(messages.reciever_user_uuid=? OR messages.sender_user_uuid=?) and  
text <> ? and 
messages.deleted_at IS NULL and
users.banned=false`,
			"private",
			user.Uuid,
			user.Uuid,
			"",
		).
		Count(&count)
	return count
}

func CountUndreadPrivateMessages(user User) int {
	var (
		count int
	)
	database.
		Table("messages").
		Joins("join users on messages.sender_user_uuid=users.uuid").
		Where(
			`messages.type=? and
	messages.reciever_user_uuid=? and
	is_read_by_reciever=? and
	text <> ? and
	messages.deleted_at IS NULL and
	users.banned=false`,
			"private", user.Uuid, 0, "",
		).
		Count(&count)
	return count
}

func CountUndreadSupportMessages(user User) int {
	var (
		count int
	)
	database.
		Table("messages").
		Where(
			"messages.type=? and (messages.reciever_user_uuid=?) and is_read_by_reciever=? and text <> ? and deleted_at IS NULL",
			"support",
			user.Uuid,
			0,
			"",
		).
		Count(&count)
	return count
}

func FindSellerThreads() Threads {
	threads := FindThreadsByType("store")
	return threads
}

func CreateThread(
	threadType, uuid, title, text string,
	senderUser *User, recieverUser *User, save bool,
) (*Thread, error) {

	if uuid != "" {
		thread, _ := FindThreadByUuid(uuid)
		if thread != nil {
			return nil, errors.New("Thread already exists")
		}
	} else {
		uuid = threadType + "-" + util.GenerateUuid()
	}

	thread := &Thread{Message: Message{
		Uuid:               uuid,
		Text:               text,
		Title:              title,
		CreatedAtTimestamp: time.Now(),

		Type: threadType,
	}}
	if senderUser != nil {
		thread.SenderUserUuid = senderUser.Uuid
	}
	if recieverUser != nil {
		thread.RecieverUserUuid = recieverUser.Uuid
	}

	err := thread.Validate()
	if err != nil {
		return thread, err
	}

	if save {
		err = thread.Save()
		if err != nil {
			return thread, err
		}
	}

	return thread, err
}

func CreateMessage(text string, thread Thread, user User) (*Message, error) {

	message := Message{
		Uuid:               thread.Uuid + "-" + util.GenerateUuid(),
		ParentUuid:         thread.Uuid,
		Text:               text,
		Type:               thread.Type,
		RecieverUserUuid:   thread.RecieverUserUuid,
		SenderUserUuid:     user.Uuid,
		SenderUser:         user,
		CreatedAtTimestamp: time.Now(),
	}

	err := message.Validate()
	if err != nil {
		return &message, err
	}

	err = message.Save()
	if err != nil {
		return &message, err
	}

	return &message, nil
}

func FindMessageByUuid(uuid string) (*Message, error) {
	var (
		message Message
	)
	err := database.
		Where("uuid=?", uuid).
		Preload("SenderUser").
		Preload("RecieverUser").
		Find(&message).Error
	if err != nil {
		return nil, err
	}

	return &message, err
}

func FindThreadByUuid(uuid string) (*Thread, error) {
	message, err := FindMessageByUuid(uuid)
	if err != nil {
		return nil, err
	}
	return &Thread{Message: *message}, nil
}

func FindMessageboardSectionByID(id int) (*MessageboardSection, error) {
	var (
		section MessageboardSection
	)

	err := database.
		Where("id = ?", id).
		Find(&section).Error

	if err != nil {
		return nil, err
	}

	return &section, err
}

func FindAllMessageboardSections() []MessageboardSection {
	var sections []MessageboardSection
	database.
		Find(&sections)
	return sections
}

/*
	Factory Quieries
*/

func GetStoreThread(user Seller) (*Thread, error) {
	threadUuid := "store-" + user.Uuid
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "store" {
		return thread, nil
	}
	return CreateThread(
		"store",
		threadUuid,
		"Store thread @"+user.Username,
		"",
		user.User,
		nil,
		true,
	)
}

func GetDisputeClaimThread(disputeClaim DisputeClaim) (*Thread, error) {
	threadUuid := fmt.Sprintf("dispute-%s-%d", disputeClaim.DisputeUuid, disputeClaim.ID)
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "dispute" {
		return thread, nil
	}
	return CreateThread(
		"dispute",
		threadUuid,
		threadUuid,
		"",
		nil,
		nil,
		true,
	)
}

func GetShoutboxThread(lang string) (*Thread, error) {
	if lang == "" {
		lang = "en"
	}
	threadUuid := "shoutbox-" + lang
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "shoutbox" {
		return thread, nil
	}
	l := GetLocalization(lang)
	return CreateThread(
		"shoutbox",
		threadUuid,
		l.LeftMenu.Shoutbox,
		"",
		nil,
		nil,
		true,
	)

}

func GetNewsThread(lang string) (*Thread, error) {
	if lang == "" {
		lang = "en"
	}
	threadUuid := "news-" + lang
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "news" {
		return thread, nil
	}
	return CreateThread(
		"news",
		threadUuid,
		"news",
		"",
		nil,
		nil,
		true,
	)

}

func GetMessageboardThread(uuid string) (*Thread, error) {
	thread, err := FindThreadByUuid(uuid)
	if err != nil {
		return nil, err
	}
	if thread != nil && thread.Type == "messageboard" {
		return thread, nil
	}
	return nil, errors.New("No such thread or thread of wrong type")
}

func GetStaffMessageboardThread(uuid string) (*Thread, error) {
	thread, err := FindThreadByUuid(uuid)
	if err != nil {
		return nil, err
	}
	if thread != nil && thread.Type == "staff_messageboard" {
		return thread, nil
	}
	return nil, errors.New("No such thread or thread of wrong type")
}

func GetTransactionThread(transaction Transaction, message string) (*Thread, error) {
	threadUuid := "transaction-" + transaction.Uuid
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "transaction" {
		return thread, nil
	}
	return CreateThread(
		"transaction",
		threadUuid,
		"Transaction thread @"+transaction.Uuid,
		message,
		&transaction.Buyer,
		&transaction.Seller,
		true,
	)
}

func GetPrivateThread(sender, reciever User, message string, createIfNoExists bool) (*Thread, error) {

	thread := FindPrivateThread(sender, reciever)
	if thread.Uuid != "" {
		return &thread, nil
	}

	return CreateThread(
		"private",
		"private-"+util.GenerateUuid(),
		fmt.Sprintf("Private thread @%s - @%s", sender.Username, reciever.Username),
		message,
		&sender,
		&reciever,
		createIfNoExists,
	)
}

func GetVendorVerificationThread(user User, createIfNotExists bool) (*Thread, error) {
	threadUuid := "store-verification-" + user.Uuid
	thread, _ := FindThreadByUuid(threadUuid)
	if thread != nil && thread.Type == "store-verification" {
		return thread, nil
	}
	return CreateThread(
		"store-verification",
		threadUuid,
		"Vendor Verification Thread @"+user.Username,
		"",
		&user,
		nil,
		createIfNotExists,
	)
}

/*
	View models
*/

type ViewMessage struct {
	*Message
	CreatedAtStr string
	TextHTML     template.HTML
	ShortText    string
}

type ViewThread struct {
	*ViewMessage
	LastMessage            *ViewMessage
	Messages               []ViewMessage
	LastUpdatedAtStr       string
	NumberOfMessages       int
	NumberOfUnreadMessages int
	IsRead                 bool
	Pages                  []int
	TitleUser              User
	LastMessageByTitleUser bool
	IsFixProblem           bool
}

func (m Message) ViewMessage(lang string) ViewMessage {
	vm := ViewMessage{
		Message:      &m,
		CreatedAtStr: humanize.Time(m.CreatedAtTimestamp),
		TextHTML:     template.HTML(messageboardHtmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(m.Text))))),
	}

	if len(m.Text) > 30 {
		vm.ShortText = m.Text[0:30]
	} else {
		vm.ShortText = m.Text
	}

	if strings.Contains(vm.ShortText, "\n") {
		parts := strings.Split(vm.ShortText, "\n")
		for _, p := range parts {
			if p != "" {
				vm.ShortText = p
				break
			}
		}
	}

	// -----BEGIN PGP MESSAGE-----
	if strings.Contains(m.Text, "-----BEGIN PGP MESSAGE-----") ||
		strings.Contains(m.Text, "-----BEGIN PGP SIGNED MESSAGE-----") ||
		strings.Contains(m.Text, "-----BEGIN PGP PUBLIC KEY BLOCK-----") {
		vm.IsEncrypted = true
	}

	if lang == "ru" {
		vm.CreatedAtStr = util.HumanizeTimeRU(m.CreatedAtTimestamp)
	}

	return vm
}

func (t Thread) ViewThread(lang string, reciever *User) ViewThread {

	viewMessage := t.Message.ViewMessage(lang)
	order := "ASC"
	if t.Type == "shoutbox" || t.Type == "news" {
		order = "DESC"
	}
	viewMessages := t.Messages(order).ViewMessages(lang)

	if (t.Type == "support" || t.Type == "store") && len(viewMessages) > 0 {
		viewMessages = viewMessages[1:len(viewMessages)]
	}

	viewThread := ViewThread{
		ViewMessage:      &viewMessage,
		Messages:         viewMessages,
		NumberOfMessages: len(viewMessages),
		LastUpdatedAtStr: humanize.Time(t.LastUpdated),
		IsRead:           false,
	}

	viewThread.LastUpdatedAtStr = util.HumanizeTime(t.LastUpdated, lang)

	if len(viewMessages) > 0 {
		viewThread.LastMessage = &viewMessages[len(viewMessages)-1]
	}

	if reciever != nil {

		if reciever.Uuid == t.RecieverUser.Uuid {
			viewThread.TitleUser = t.SenderUser
		} else {
			viewThread.TitleUser = t.RecieverUser
		}

		if viewThread.LastMessage != nil {
			if viewThread.LastMessage.SenderUser.Uuid == viewThread.TitleUser.Uuid {
				viewThread.LastMessageByTitleUser = false
			} else {
				viewThread.LastMessageByTitleUser = true
			}
			tps, _ := FindThreadPerusalStatus(t.Uuid, reciever.Uuid)
			if tps != nil && tps.LastReadDate.After(viewThread.LastMessage.CreatedAtTimestamp) {
				viewThread.IsRead = true
			} else {
				viewThread.IsRead = false
			}

			viewThread.IsFixProblem = FindThreadSupportStatus(viewThread.LastMessage.Uuid)
		}

		if reciever.Username == "" {
			viewThread.IsRead = false
		}
	}

	numberOfPages := int(math.Ceil(float64(viewThread.NumberOfMessages) / 10.0))
	for i := 0; i < numberOfPages; i++ {
		viewThread.Pages = append(viewThread.Pages, i+1)
	}

	return viewThread
}

func FindThreadSupportStatus(messageuuid string) bool {
	var tss ThreadSupportStatus

	err := database.
		Table("thread_support_statuses").
		Where("message_uuid_mark = ?", messageuuid).
		Limit(1).
		Find(&tss).
		Error
	if err != nil {

	}
	return tss.IsFixProblem
}

func (ts Threads) ViewThreads(lang string, reciever *User) []ViewThread {
	viewThreads := []ViewThread{}

	for i, thread := range ts {
		if thread.Type == "messageboard" || (len(thread.Messages("ASC")) > 1 && thread.Type != "messageboard") {
			vt := ts[i].ViewThread(lang, reciever)
			viewThreads = append(viewThreads, vt)
		}
	}

	for i, j := 0, len(viewThreads)-1; i < j; i, j = i+1, j-1 {
		viewThreads[i], viewThreads[j] = viewThreads[j], viewThreads[i]
	}

	return viewThreads
}

func (ms Messages) ViewMessages(lang string) []ViewMessage {
	viewMessages := []ViewMessage{}
	for i, _ := range ms {
		vm := ms[i].ViewMessage(lang)
		viewMessages = append(viewMessages, vm)
	}
	return viewMessages
}

func setupThreadsViews() {
	database.Exec("DROP VIEW IF EXISTS v_threads CASCADE;")
	database.Exec(`
		CREATE VIEW v_threads AS (
			WITH thread_messages as (
				SELECT parent_uuid, MAX(created_at_timestamp) last_updated, count(*) as number_of_messages
				FROM messages
				WHERE parent_uuid <> ''
				AND (deleted_at IS NULL OR deleted_at <= '0001-01-02') 
				GROUP BY parent_uuid
				ORDER BY number_of_messages DESC
			),
			extended_thread_messages as (
				SELECT thread_messages.*, messages.uuid as last_message_uuid 
				FROM thread_messages 
				JOIN messages ON messages.parent_uuid=thread_messages.parent_uuid AND messages.created_at_timestamp=thread_messages.last_updated
			)

			SELECT  * FROM (
				SELECT
					messages.*,
					tm.last_updated,
					tm.number_of_messages,
					tm.last_message_uuid,
					u.is_admin,
					u.premium_plus,
					u.premium
				FROM	
					messages
				JOIN
					extended_thread_messages tm on tm.parent_uuid=messages.uuid
				JOIN 
					users u on u.uuid=messages.sender_user_uuid
				WHERE
					(messages.deleted_at IS NULL OR messages.deleted_at <= '0001-01-02') AND
					u.banned=false 
				UNION (
					SELECT 
						messages.*,
						created_at_timestamp AS last_updated,
						1,
						messages.uuid,
						u.is_admin,
						u.premium_plus,
						u.premium
					FROM 
						messages
					LEFT JOIN
						thread_messages tm on messages.uuid=tm.parent_uuid
					JOIN 
						users u on u.uuid=messages.sender_user_uuid
					WHERE
						(messages.deleted_at IS NULL OR messages.deleted_at <= '0001-01-02') AND
						u.banned=false AND
						messages.parent_uuid = '' AND
						(tm.parent_uuid IS NULL OR tm.parent_uuid='')
				)
			) threads ORDER BY is_admin, premium_plus, premium, created_at_timestamp ASC
	);`)

	database.Exec("DROP VIEW IF EXISTS v_thread_counts CASCADE;")
	database.Exec(`
		CREATE VIEW v_thread_counts AS (
			SELECT 
				parent_uuid, count(*) as number_of_messages 
			FROM 
				messages 
			GROUP BY
				parent_uuid
	);`)
}
