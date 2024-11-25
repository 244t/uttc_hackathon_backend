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
	RegisterUser(user model.Profile) error
	GetUserProfile(userId string) (model.Profile,error)
	GetFollowing(userId string) ([]model.Profile,error)
	GetFollowers(userId string)([]model.Profile,error)
	Follow(follow model.Follow) error
	UnFollow(unfollow model.UnFollow) error
	UpdateProfile(user model.Profile) error
	GetUserPosts(userId string) ([]model.Post,error)
}


//TweetDAOのインスタンスを返す
func NewTweetDAO (db *sql.DB) *TweetDAO{
	return &TweetDAO{DB:db}
}


func (dao *TweetDAO) RegisterUser(user model.Profile) error{
	_ ,err := dao.DB.Exec("INSERT INTO user (user_id, name, bio,profile_img_url,header_img_url,location) VALUES (?, ?, ?,?,?,?)", user.Id, user.Name, user.Bio,user.ImgUrl,user.HeaderUrl,user.Location)
	return err
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
		SELECT u.user_id, u.name, u.bio, u.profile_img_url
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
		SELECT u.user_id, u.name, u.bio, u.profile_img_url
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

func (dao *TweetDAO) Follow(follow model.Follow) error{
	now := time.Now()
	_ ,err := dao.DB.Exec("INSERT INTO follower (user_id, following_user_id,created_at) VALUES (?, ?, ?)", follow.UserId,follow.FollowingId,now)
	return err
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
		user.Name, user.Bio, user.ImgUrl, user.Id)
	return err
}

func (dao *TweetDAO) GetUserPosts(userId string) ([]model.Post,error){
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