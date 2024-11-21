package dao

import(
	"database/sql"
	"myproject/model"
	"log"
	"time"
	"fmt"
)

type PostDAO struct {
	DB *sql.DB
}


type PostDAOInterface interface{
	CreatePost(newpost model.Post) error
	GetPost(postId string) (model.Post,error)
	UpdatePost(updatepost model.Update) error
	DeletePost(deletePost model.Delete) error
	LikePost(l model.Like) error
	UnLikePost(u model.Like) error
	GetLikesPost(postId string) (model.Likes,error)
	Timeline(userId string) ([]model.Post,error)
}

func NewPostDAO (db *sql.DB) *PostDAO{
	return &PostDAO{DB:db}
}

func (dao *PostDAO) CreatePost(newpost model.Post) error {
	_ ,err := dao.DB.Exec("INSERT INTO post (post_id, user_id, content,img_url,created_at,edited_at,deleted_at,parent_post_id) VALUES (?, ?, ?,?, ?, ?, ?, ?)", newpost.PostId,newpost.UserId,newpost.Content,newpost.ImgUrl,newpost.CreatedAt,newpost.EditedAt,newpost.DeletedAt,newpost.ParentPostId)
	return err
}

func (dao *PostDAO) GetPost(postId string) (model.Post, error){
	var post model.Post
	var createdAtRaw []byte 
	var editedAtRaw []byte
	var deletedAtRaw []byte
	query := `
	SELECT post_id, user_id, content,img_url, created_at, edited_at, deleted_at, parent_post_id 
	FROM post 
	WHERE post_id = ?
	`
	err := dao.DB.QueryRow(query, postId).Scan(
        &post.PostId,
        &post.UserId,
        &post.Content,
		&post.ImgUrl,
        &createdAtRaw,  // created_atを[]byteで受け取る
		&editedAtRaw,   // edited_atを[]byteで受け取る
		&deletedAtRaw,  // deleted_atを[]byteで受け取る
        &post.ParentPostId,
    )
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからなかった場合
			return model.Post{}, nil  // 空の構造体を返す
		}

		// その他のエラー
		log.Printf("Error fetching post for postId %s: %v", postId, err)
		return model.Post{}, fmt.Errorf("could not fetch post: %w", err)  // ラップしたエラーを返す
	}
	// createdAtRaw を time.Time に変換
	// "2006-01-02 15:04:05" は一般的なDATETIMEのフォーマットです
	createdAtStr := string(createdAtRaw)
	post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		log.Printf("Error parsing created_at: %v", err)
		return model.Post{}, fmt.Errorf("could not parse created_at: %w", err)
	}
	// editedAtRaw を sql.NullTime に変換
	if len(editedAtRaw) > 0 {
		editedAtStr := string(editedAtRaw)
		editedAt, err := time.Parse("2006-01-02 15:04:05", editedAtStr)
		if err != nil {
			log.Printf("Error parsing edited_at: %v", err)
			return model.Post{}, fmt.Errorf("could not parse edited_at: %w", err)
		}
		post.EditedAt = sql.NullTime{Time: editedAt, Valid: true}
	} else {
		post.EditedAt = sql.NullTime{Valid: false}
	}

	// deletedAtRaw を sql.NullTime に変換
	if len(deletedAtRaw) > 0 {
		deletedAtStr := string(deletedAtRaw)
		deletedAt, err := time.Parse("2006-01-02 15:04:05", deletedAtStr)
		if err != nil {
			log.Printf("Error parsing deleted_at: %v", err)
			return model.Post{}, fmt.Errorf("could not parse deleted_at: %w", err)
		}
		post.DeletedAt = sql.NullTime{Time: deletedAt, Valid: true}
	} else {
		post.DeletedAt = sql.NullTime{Valid: false}
	}

	return post, nil
}

func (dao *PostDAO) UpdatePost(updatePost model.Update) error {
	// SQL クエリを作成
	query := `
		UPDATE post
		SET content = ?, img_url = ?, edited_at = ?
		WHERE post_id = ?
	`
	// SQL クエリを実行
	_, err := dao.DB.Exec(query, updatePost.Content, updatePost.ImgUrl,updatePost.EditedAt, updatePost.PostId)
	if err != nil {
		// エラーが発生した場合はそのエラーを返す
		return fmt.Errorf("could not update post: %w", err)
	}

	return err
}

func (dao *PostDAO) DeletePost(deletePost model.Delete)error{
	query := `
		UPDATE post
		SET deleted_at = ?
		WHERE post_id = ?
	`
	_, err := dao.DB.Exec(query,deletePost.DeletedAt,deletePost.PostId)
	if err != nil {
		// エラーが発生した場合はそのエラーを返す
		return fmt.Errorf("could not delete post: %w", err)
	}

	return err
}

func (dao *PostDAO) LikePost(l model.Like) error{
	_ ,err := dao.DB.Exec("INSERT INTO `like` (post_id, user_id,created_at) VALUES (?, ?, ?)", l.PostId,l.UserId,l.CreatedAt)
	return err
}

func (dao *PostDAO) UnLikePost(u model.Like) error{
	_, err := dao.DB.Exec("DELETE FROM `like` WHERE post_id = ? AND user_id = ?", u.PostId, u.UserId)
    return err 
}

func (dao *PostDAO) GetLikesPost (postId string) (model.Likes,error) {
	var l model.Likes
	query := "SELECT COUNT(*) FROM `like` WHERE post_id = ?"
	err := dao.DB.QueryRow(query, postId).Scan(&l.Likes)
	if err != nil {
		fmt.Errorf("error executing query: %w", err)
        return model.Likes{}, err
    }
	return l,nil
}

func (dao *PostDAO) Timeline(userId string) ([]model.Post, error) {
    // 1. フォローしているユーザーの投稿を取得
    query := `
        SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.edited_at, p.deleted_at, p.parent_post_id
        FROM post p
        JOIN follower f ON f.following_user_id = p.user_id
        WHERE f.user_id = ?
        `
    
    // クエリ実行
    rows, err := dao.DB.Query(query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    // 結果を格納するスライス
    var posts []model.Post
    for rows.Next() {
		var post model.Post
		var createdAtRaw []byte
		var editedAtRaw []byte
		var deletedAtRaw []byte
		err := rows.Scan(
			&post.PostId,
			&post.UserId,
			&post.Content,
			&post.ImgUrl,
			&createdAtRaw,  // created_atを[]byteで受け取る
			&editedAtRaw,   // edited_atを[]byteで受け取る
			&deletedAtRaw,  // deleted_atを[]byteで受け取る
			&post.ParentPostId,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				// ユーザーが見つからなかった場合
				return nil, err  // 空の構造体を返す
			}
			// その他のエラー
			log.Printf("Error fetching timeline: %v", err)
			return nil, err
		}

		// createdAtRaw を time.Time に変換
		if len(createdAtRaw) > 0 {
			createdAtStr := string(createdAtRaw)
			post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				log.Printf("Error parsing created_at: %v", err)
				return nil, err
			}
		} else {
			log.Printf("created_at is NULL")
			return nil,err
		}

		// editedAtRaw を sql.NullTime に変換
		if len(editedAtRaw) > 0 {
			editedAtStr := string(editedAtRaw)
			editedAt, err := time.Parse("2006-01-02 15:04:05", editedAtStr)
			if err != nil {
				log.Printf("Error parsing edited_at: %v", err)
				return nil,err
			}
			post.EditedAt = sql.NullTime{Time: editedAt, Valid: true}
		} else {
			post.EditedAt = sql.NullTime{Valid: false}  // NULL 値
		}

		// deletedAtRaw を sql.NullTime に変換
		if len(deletedAtRaw) > 0 {
			deletedAtStr := string(deletedAtRaw)
			deletedAt, err := time.Parse("2006-01-02 15:04:05", deletedAtStr)
			if err != nil {
				log.Printf("Error parsing deleted_at: %v", err)
				return nil,err
			}
			post.DeletedAt = sql.NullTime{Time: deletedAt, Valid: true}
		} else {
			post.DeletedAt = sql.NullTime{Valid: false}  // NULL 値
		}
        posts = append(posts, post)
    }

    return posts, nil
}