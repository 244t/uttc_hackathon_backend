package model
import(
	"time"
	"database/sql"
)

// type Post struct {
// 	PostId string `json:"post_id"`
// 	UserId string `json:"user_id"`
// 	Content string `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// 	EditedAt *time.Time `json:"edited_at"`
// 	DeletedAt *time.Time `json:"deleted_at"`
// 	ParentPostId *string `json:"parent_post_id"`
// }

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

type Likes struct {
	Likes []string `json:"user_ids"`
}

