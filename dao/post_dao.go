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
	GetChildPostCount(postId string) (int,error)
	CreatePost(newpost model.Post) error
	GetPost(postId string) (model.PostWithReplyCounts,error)
	UpdatePost(updatepost model.Update) error
	DeletePost(deletePost model.Delete) error
	GetReplyPost(postId string)([]model.PostWithReplyCounts,error)
	LikePost(l model.Like, n model.Notification) error
	UnLikePost(u model.Like) error
	GetLikesPost(postId string) ([]string,error)
	Timeline(userId string) ([]model.PostWithReplyCounts,error)
}

func NewPostDAO (db *sql.DB) *PostDAO{
	return &PostDAO{DB:db}
}

func (dao *PostDAO) GetChildPostCount(postId string) (int, error) {
	// 子ポストの数を取得するクエリ
	query := `
	SELECT COUNT(*) 
	FROM post 
	WHERE parent_post_id = ?
	`

	var count int
	err := dao.DB.QueryRow(query, postId).Scan(&count)
	if err != nil {
		log.Printf("Error counting child posts for postId %s: %v", postId, err)
		return 0, fmt.Errorf("could not count child posts: %w", err)
	}
	return count, nil
}


func (dao *PostDAO) CreatePost(newpost model.Post) error {
	_ ,err := dao.DB.Exec("INSERT INTO post (post_id, user_id, content,img_url,created_at,edited_at,deleted_at,parent_post_id) VALUES (?, ?, ?,?, ?, ?, ?, ?)", newpost.PostId,newpost.UserId,newpost.Content,newpost.ImgUrl,newpost.CreatedAt,newpost.EditedAt,newpost.DeletedAt,newpost.ParentPostId)
	return err
}

func (dao *PostDAO) GetPost(postId string) (model.PostWithReplyCounts, error) {
	// ポストを格納する変数
	var post model.PostWithReplyCounts
	var createdAtRaw []byte
	var editedAtRaw []byte
	var deletedAtRaw []byte
	var childPostCount int

	// ポストを取得するクエリ
	query := `
	SELECT post_id, user_id, content, img_url, created_at, edited_at, deleted_at, parent_post_id 
	FROM post 
	WHERE post_id = ?
	`
	// クエリ実行
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
			return model.PostWithReplyCounts{}, nil  // 空の構造体を返す
		}

		// その他のエラー
		log.Printf("Error fetching post for postId %s: %v", postId, err)
		return model.PostWithReplyCounts{}, fmt.Errorf("could not fetch post: %w", err)  // ラップしたエラーを返す
	}

	// createdAtRaw を time.Time に変換
	createdAtStr := string(createdAtRaw)
	post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		log.Printf("Error parsing created_at: %v", err)
		return model.PostWithReplyCounts{}, fmt.Errorf("could not parse created_at: %w", err)
	}

	// editedAtRaw を sql.NullTime に変換
	if len(editedAtRaw) > 0 {
		editedAtStr := string(editedAtRaw)
		editedAt, err := time.Parse("2006-01-02 15:04:05", editedAtStr)
		if err != nil {
			log.Printf("Error parsing edited_at: %v", err)
			return model.PostWithReplyCounts{}, fmt.Errorf("could not parse edited_at: %w", err)
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
			return model.PostWithReplyCounts{}, fmt.Errorf("could not parse deleted_at: %w", err)
		}
		post.DeletedAt = sql.NullTime{Time: deletedAt, Valid: true}
	} else {
		post.DeletedAt = sql.NullTime{Valid: false}
	}

	// 2. 子ポストの数を取得（GetChildPostCountを使用）
	childPostCount, err = dao.GetChildPostCount(postId)
	if err != nil {
		log.Printf("Error counting child posts for postId %s: %v", postId, err)
		return model.PostWithReplyCounts{}, fmt.Errorf("could not count child posts: %w", err)
	}

	post.ReplyCounts = childPostCount


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

// func (dao *PostDAO) LikePost(l model.Like, n model.Notification) error {
//     // Insert like into the database
//     _, err := dao.DB.Exec("INSERT INTO `like` (post_id, user_id, created_at) VALUES (?, ?, ?)", l.PostId, l.UserId, l.CreatedAt)
//     if err != nil {
//         return fmt.Errorf("could not like post: %w", err)
//     }

//     // Insert notification into the database
//     _, err = dao.DB.Exec("INSERT INTO notification (notification_id, user_id, flag, action_user_id) VALUES (?, ?, ?, ?)", n.NotificationId, n.UserId, n.Flag, n.ActionUserId)
//     if err != nil {
//         return fmt.Errorf("could not create notification: %w", err)
//     }

//     return nil
// }
func (dao *PostDAO) LikePost(l model.Like, n model.Notification) error {
    // Insert like into the database
    _, err := dao.DB.Exec("INSERT INTO `like` (post_id, user_id, created_at) VALUES (?, ?, ?)", l.PostId, l.UserId, l.CreatedAt)
    if err != nil {
        return fmt.Errorf("could not like post: %w", err)
    }

    // Retrieve the user_id from the post table based on post_id
    var postUserId string
    err = dao.DB.QueryRow("SELECT user_id FROM post WHERE post_id = ?", l.PostId).Scan(&postUserId)
    if err != nil {
        return fmt.Errorf("could not retrieve user_id from post: %w", err)
    }

    // Update n.UserId with the retrieved user_id from the post
    n.UserId = postUserId

    // Insert notification into the database
    _, err = dao.DB.Exec("INSERT INTO notification (notification_id, user_id, flag, action_user_id) VALUES (?, ?, ?, ?)", n.NotificationId, n.UserId, n.Flag, n.ActionUserId)
    if err != nil {
        return fmt.Errorf("could not create notification: %w", err)
    }

    return nil
}



func (dao *PostDAO) UnLikePost(u model.Like) error{
	_, err := dao.DB.Exec("DELETE FROM `like` WHERE post_id = ? AND user_id = ?", u.PostId, u.UserId)
    return err 
}

func (dao *PostDAO) GetLikesPost(postId string) ([]string, error) {
    // user_idを格納するスライス
    var userIds []string
    
    // `user_id` を取得するクエリ
    query := "SELECT user_id FROM `like` WHERE post_id = ?"
    
    // クエリを実行し、取得したuser_idをスライスに追加
    rows, err := dao.DB.Query(query, postId)
    if err != nil {
        fmt.Errorf("error executing query: %w", err)
        return nil, err
    }
    defer rows.Close()
    
    // 取得したuser_idをスライスに格納
    for rows.Next() {
        var userId string
        if err := rows.Scan(&userId); err != nil {
            fmt.Errorf("error scanning row: %w", err)
            return nil, err
        }
        userIds = append(userIds, userId)
    }
    
    // rowsの走査中にエラーが発生していないか確認
    if err := rows.Err(); err != nil {
        fmt.Errorf("error during row iteration: %w", err)
        return nil, err
    }

    // userIdsスライスを返す
    return userIds, nil
}


func (dao *PostDAO) Timeline(userId string) ([]model.PostWithReplyCounts, error) {
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
	var posts []model.PostWithReplyCounts
	for rows.Next() {
		var post model.PostWithReplyCounts
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
			return nil, err
		}

		// editedAtRaw を sql.NullTime に変換
		if len(editedAtRaw) > 0 {
			editedAtStr := string(editedAtRaw)
			editedAt, err := time.Parse("2006-01-02 15:04:05", editedAtStr)
			if err != nil {
				log.Printf("Error parsing edited_at: %v", err)
				return nil, err
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
				return nil, err
			}
			post.DeletedAt = sql.NullTime{Time: deletedAt, Valid: true}
		} else {
			post.DeletedAt = sql.NullTime{Valid: false}  // NULL 値
		}

		// 2. 子ポスト数を取得
		childPostCount, err := dao.GetChildPostCount(post.PostId)
		if err != nil {
			return nil, fmt.Errorf("could not fetch child post count: %w", err)
		}

		// 子ポスト数をポスト構造体に設定
		post.ReplyCounts = childPostCount

		// 結果を追加
		posts = append(posts, post)
	}

	// タイムラインを返す
	return posts, nil
}

func (dao *PostDAO) GetReplyPost(postId string) ([]model.PostWithReplyCounts, error) {
    // 1. 指定されたpostIdをparent_post_idとして持つ投稿を取得
    query := `
        SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.edited_at, p.deleted_at, p.parent_post_id
        FROM post p
        WHERE p.parent_post_id = ?
    `
    
    // クエリ実行
    rows, err := dao.DB.Query(query, postId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // 結果を格納するスライス
    var posts []model.PostWithReplyCounts
    for rows.Next() {
        var post model.PostWithReplyCounts
        var createdAtRaw []byte
        var editedAtRaw []byte
        var deletedAtRaw []byte
        var childPostCount int
        
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
                // 子投稿が見つからなかった場合
                return nil, nil // 空のスライスを返す
            }
            log.Printf("Error fetching reply posts: %v", err)
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
            return nil, err
        }

        // editedAtRaw を sql.NullTime に変換
        if len(editedAtRaw) > 0 {
            editedAtStr := string(editedAtRaw)
            editedAt, err := time.Parse("2006-01-02 15:04:05", editedAtStr)
            if err != nil {
                log.Printf("Error parsing edited_at: %v", err)
                return nil, err
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
                return nil, err
            }
            post.DeletedAt = sql.NullTime{Time: deletedAt, Valid: true}
        } else {
            post.DeletedAt = sql.NullTime{Valid: false}  // NULL 値
        }

        // 子ポストの数を取得（GetChildPostCountを使用）
        childPostCount, err = dao.GetChildPostCount(post.PostId)
        if err != nil {
            log.Printf("Error counting child posts for postId %s: %v", post.PostId, err)
            return nil, fmt.Errorf("could not count child posts: %w", err)
        }
		post.ReplyCounts = childPostCount

        posts = append(posts, post)
    }

    return posts, nil
}
