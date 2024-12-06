package dao

import(
	"database/sql"
	"fmt"
	"time"
	"myproject/model"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


type TweetDAO struct {
	DB *sql.DB
}

type TweetDAOInterface interface{
	GetChildPostCount(postId string) (int,error)
	GetUserProfile(userId string) (model.Profile,error)
	GetFollowing(userId string) ([]model.Profile,error)
	GetFollowers(userId string)([]model.Profile,error)
	Follow(follow model.Follow,n model.Notification) error
	UnFollow(unfollow model.UnFollow) error
	UpdateProfile(user model.Profile) error
	GetUserPosts(userId string) ([]model.PostWithReplyCounts,error)
	SearchUser(searchWord string)([]model.Profile,error)
	Notification(userId string) ([]model.NotificationInfo,error)
}


//TweetDAOのインスタンスを返す
func NewTweetDAO (db *sql.DB) *TweetDAO{
	return &TweetDAO{DB:db}
}

func (dao *TweetDAO) GetChildPostCount(postId string) (int, error) {
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


//user_idをもとにユーザープロフィールを得る
func (dao *TweetDAO) GetUserProfile(userId string) (model.Profile, error) {
	var prof model.Profile
	err := dao.DB.QueryRow("SELECT user_id, name, bio, profile_img_url,header_img_url,location FROM user WHERE user_id = ?", userId).Scan(&prof.Id, &prof.Name, &prof.Bio,&prof.ImgUrl,&prof.HeaderUrl,&prof.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからなかった場合
			return model.Profile{}, nil  // 空の構造体を返す
		}

		// その他のエラー
		log.Printf("Error fetching user profile for userId %s: %v", userId, err)
		return model.Profile{}, fmt.Errorf("could not fetch user profile: %w", err)  // ラップしたエラーを返す
	}

	return prof, nil
}

func (dao *TweetDAO) GetFollowing(userId string) ([]model.Profile, error) {
	// INNER JOINを使用して、follower テーブルと user テーブルを結合
	query := `
		SELECT u.user_id, u.name, u.bio, u.profile_img_url, u.header_img_url, location
		FROM follower f
		INNER JOIN user u ON f.following_user_id = u.user_id
		WHERE f.user_id = ?
	`

	// フォローしているユーザーのプロフィールを格納するスライス
	var profiles []model.Profile

	// フォローしているユーザーのプロフィール情報を取得
	rows, err := dao.DB.Query(query, userId)
	if err != nil {
		log.Printf("Error fetching following user profiles for userId %s: %v", userId, err)
		return nil, fmt.Errorf("could not fetch following user profiles: %w", err)
	}
	defer rows.Close()

	// rowsからプロフィール情報を読み取って、profilesスライスに追加
	for rows.Next() {
		var profile model.Profile
		if err := rows.Scan(&profile.Id, &profile.Name, &profile.Bio,&profile.ImgUrl,&profile.HeaderUrl,&profile.Location); err != nil {
			log.Printf("Error scanning profile: %v", err)
			return nil, fmt.Errorf("could not scan profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	// フォローしているユーザーのプロフィールを返す
	return profiles, nil
}


func (dao *TweetDAO) GetFollowers(userId string) ([]model.Profile, error) {
	// フォローしているユーザーのプロフィールを取得するためのSQLクエリ
	query := `
		SELECT u.user_id, u.name, u.bio, u.profile_img_url, u.header_img_url, u.location
		FROM user u
		INNER JOIN follower f ON u.user_id = f.user_id
		WHERE f.following_user_id = ?
	`

	// フォローしているユーザーのプロフィールを格納するスライス
	var profiles []model.Profile

	// データベースから情報を取得
	rows, err := dao.DB.Query(query, userId)
	if err != nil {
		log.Printf("Error fetching followers for userId %s: %v", userId, err)
		return nil, fmt.Errorf("could not fetch followers: %w", err)
	}
	defer rows.Close()

	// 取得した各行を処理
	for rows.Next() {
		var profile model.Profile
		if err := rows.Scan(&profile.Id, &profile.Name, &profile.Bio,&profile.ImgUrl,&profile.HeaderUrl,&profile.Location); err != nil {
			log.Printf("Error scanning profile for userId %s: %v", userId, err)
			continue
		}
		// プロフィールをスライスに追加
		profiles = append(profiles, profile)
	}

	// エラーがあれば返す
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// フォローしているユーザーのプロフィールを返す
	return profiles, nil
}

func (dao *TweetDAO) Follow(follow model.Follow,n model.Notification) error{
	now := time.Now()
	_ ,err := dao.DB.Exec("INSERT INTO follower (user_id, following_user_id,created_at) VALUES (?, ?, ?)", follow.UserId,follow.FollowingId,now)
	if err != nil {
		return fmt.Errorf("could not register follow notification")
	}
	_, err = dao.DB.Exec("INSERT INTO notification (notification_id, user_id, flag, action_user_id) VALUES (?, ?, ?, ?)", n.NotificationId, n.UserId, n.Flag, n.ActionUserId)
    if err != nil {
        return fmt.Errorf("could not register follow notification: %w", err)
    }
	return nil
}

func (dao *TweetDAO) UnFollow(unfollow model.UnFollow) error {
	query := `
		DELETE FROM follower
		WHERE user_id = ? AND following_user_id = ?
	`
	// クエリの実行
	_, err := dao.DB.Exec(query, unfollow.UserId, unfollow.FollowingId)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		return fmt.Errorf("could not unfollow user: %v", err)
	}

	return nil
}

func (dao *TweetDAO) UpdateProfile(user model.Profile) error {
	_, err := dao.DB.Exec(`
		UPDATE user
		SET name = ?, bio = ?, profile_img_url = ?, header_img_url = ?, location = ?
		WHERE user_id = ?`,
		user.Name, user.Bio, user.ImgUrl,user.HeaderUrl,user.Location,user.Id)
	return err
}

func (dao *TweetDAO) GetUserPosts(userId string) ([]model.PostWithReplyCounts,error){
	// 1. フォローしているユーザーの投稿を取得
	query := `
		SELECT post_id, user_id, content, img_url, created_at, edited_at, deleted_at, parent_post_id
		FROM post 
		WHERE user_id = ?
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

		// 2. 子ポスト数を取得
		childPostCount, err := dao.GetChildPostCount(post.PostId)
		if err != nil {
			return nil, fmt.Errorf("could not fetch child post count: %w", err)
		}

		// 子ポスト数をポスト構造体に設定
		post.ReplyCounts = childPostCount


		posts = append(posts, post)
	}

	return posts, nil
}

func (dao *TweetDAO) SearchUser(sw string) ([]model.Profile, error) {
    // 名前がswで始まるユーザーを検索するためのSQLクエリ
    query := `
        SELECT user_id, name, bio, profile_img_url, header_img_url, location
        FROM user
        WHERE name LIKE ?
    `

    // swの末尾に%を付けてLIKE検索を準備
    sw = sw + "%"

    // 結果を格納するスライス
    var profiles []model.Profile

    // データベースから情報を取得
    rows, err := dao.DB.Query(query, sw)
    if err != nil {
        log.Printf("Error fetching users with name starting with '%s': %v", sw, err)
        return nil, fmt.Errorf("could not fetch users: %w", err)
    }
    defer rows.Close()

    // 取得した各行を処理
    for rows.Next() {
        var profile model.Profile
        if err := rows.Scan(&profile.Id, &profile.Name, &profile.Bio, &profile.ImgUrl, &profile.HeaderUrl, &profile.Location); err != nil {
            log.Printf("Error scanning profile: %v", err)
            continue
        }
        // プロフィールをスライスに追加
        profiles = append(profiles, profile)
    }

    // エラーがあれば返す
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating over rows: %v", err)
        return nil, err
    }

    // 検索結果のユーザーのプロフィールを返す
    return profiles, nil
}

func (dao *TweetDAO) Notification(userId string) ([]model.NotificationInfo, error) {
	// notificationテーブルからuserIdと一致するすべての行を取得
	rows, err := dao.DB.Query(`
		SELECT n.notification_id, n.flag, n.action_user_id, u.profile_img_url, u.name
		FROM notification n
		LEFT JOIN user u ON n.action_user_id = u.user_id
		WHERE n.user_id = ?`, userId)

	if err != nil {
		return nil, fmt.Errorf("could not fetch dao notifications: %w", err)
	}
	defer rows.Close()

	// 結果を格納するスライス
	var notifications []model.NotificationInfo

	// クエリの結果を構造体にマッピング
	for rows.Next() {
		var notification model.NotificationInfo
		// action_user_idとprofile_imgも取得する
		if err := rows.Scan(&notification.NotificationId, &notification.Flag, &notification.UserId, &notification.Name,&notification.UserProfileImg); err != nil {
			return nil, fmt.Errorf("could not scan notification row: %w", err)
		}
		// notificationsスライスに追加
		notifications = append(notifications, notification)
	}

	// 結果の確認
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	// 取得した通知のリストを返す
	return notifications, nil
}
