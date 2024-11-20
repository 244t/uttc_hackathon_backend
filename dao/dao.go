package dao

import(
	"database/sql"
	"fmt"
	"myproject/model"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)


type TweetDAO struct {
	DB *sql.DB
}

type TweetDAOInterface interface{
	RegisterUser(user model.Profile) error
	GetUserProfile(userId string) (model.Profile,error)
}

//DBへの接続を初期化
func NewDBConnection()(*sql.DB,error) {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
		return nil,err
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
		return nil,err
	}
	return db,nil

	// ////localとつなげるとき
	// mysqlUser := "user"
	// mysqlUserPwd := "password"
	// mysqlDatabase := "mydatabase"
	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase))
	// if err != nil {
	// 	log.Fatalf("fail: sql.Open, %v\n", err)
	// }
	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("fail: _db.Ping, %v\n", err)
	// }
	// return db,nil

}

//TweetDAOのインスタンスを返す
func NewTweetDAO (db *sql.DB) *TweetDAO{
	return &TweetDAO{DB:db}
}


func (dao *TweetDAO) RegisterUser(user model.Profile) error{
	_ ,err := dao.DB.Exec("INSERT INTO user (user_id, name, bio,firebase_uid,profile_img_url) VALUES (?, ?, ?, ?,?)", user.Id, user.Name, user.Bio,user.FireBaseId,"")
	return err
}

//user_idをもとにユーザープロフィールを得る
func (dao *TweetDAO) GetUserProfile(userId string) (model.Profile, error) {
	var prof model.Profile
	err := dao.DB.QueryRow("SELECT user_id, name, bio FROM user WHERE user_id = ?", userId).Scan(&prof.Id, &prof.Name, &prof.Bio)
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
