package model
import(
	"time"
	"database/sql"
)

type Post struct {
	PostId       string       `json:"post_id"`
	UserId       string       `json:"user_id"`
	Content      string       `json:"content"`
	ImgUrl		string `json:"img_url"`
	CreatedAt    time.Time    `json:"created_at"`   // 必須カラム、time.Time型
	EditedAt     sql.NullTime `json:"edited_at"`    // NULL可能カラム、sql.NullTime型
	DeletedAt    sql.NullTime `json:"deleted_at"`   // NULL可能カラム、sql.NullTime型
	ParentPostId  sql.NullString     `json:"parent_post_id"`
}

type Update struct{
	PostId string `json:"post_id"`
	Content string `json:"content"`
	ImgUrl string `json:"img_url"`
	EditedAt sql.NullTime `json:"edited_at"`  
}

type Delete struct{
	PostId string `json:"post_id"`
	DeletedAt sql.NullTime `json:"deleted_at"`  
}

type Like struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	UserId string `json:"user_id"`
	NotificationId string `json:"notification_id"`
	Flag string `json:"flag"`
	ActionUserId string `"json:action_user_id"`
}

type Likes struct {
	Likes []string `json:"user_ids"`
}

type PostWithReplyCounts struct {
	PostId       string       `json:"post_id"`
	UserId       string       `json:"user_id"`
	Content      string       `json:"content"`
	ImgUrl		string `json:"img_url"`
	CreatedAt    time.Time    `json:"created_at"`   // 必須カラム、time.Time型
	EditedAt     sql.NullTime `json:"edited_at"`    // NULL可能カラム、sql.NullTime型
	DeletedAt    sql.NullTime `json:"deleted_at"`   // NULL可能カラム、sql.NullTime型
	ParentPostId  sql.NullString     `json:"parent_post_id"`
	ReplyCounts int           `json:"reply_counts"` // 子ポスト数
}