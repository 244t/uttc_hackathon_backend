// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/oklog/ulid/v2"
// )

// type UserResForHTTPGet struct {
// 	Id   string `json:"id"`
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// type UserRequest struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// // ① GoプログラムからMySQLへ接続
// var db *sql.DB

// func init() {
// 	// ①-1
// 	// DB接続のための準備
// 	mysqlUser := os.Getenv("MYSQL_USER")
// 	mysqlPwd := os.Getenv("MYSQL_PWD")
// 	mysqlHost := os.Getenv("MYSQL_HOST")
// 	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

// 	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
// 	_db, err := sql.Open("mysql", connStr)

// 	if err != nil {
// 		log.Fatalf("fail: sql.Open, %v\n", err)
// 	}
// 	// ①-3
// 	if err := _db.Ping(); err != nil {
// 		log.Fatalf("fail: _db.Ping, %v\n", err)
// 	}
// 	db = _db
// }

// // ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
// func handler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		// ②-1
// 		queryParams := r.URL.Query()
// 		name := queryParams["name"][0] // To be filled
// 		if name == "" {
// 			log.Println("fail: name is empty")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		// ②-2
// 		rows, err := db.Query("SELECT id, name, age FROM user WHERE name = ?", name)
// 		if err != nil {
// 			log.Printf("fail: db.Query, %v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// ②-3
// 		users := make([]UserResForHTTPGet, 0)
// 		for rows.Next() {
// 			var u UserResForHTTPGet
// 			if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
// 				log.Printf("fail: rows.Scan, %v\n", err)

// 				if err := rows.Close(); err != nil { // 500を返して終了するが、その前にrowsのClose処理が必要
// 					log.Printf("fail: rows.Close(), %v\n", err)
// 				}
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			users = append(users, u)
// 		}

// 		// ②-4
// 		bytes, err := json.Marshal(users)
// 		if err != nil {
// 			log.Printf("fail: json.Marshal, %v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(bytes)

// 	case http.MethodPost:
// 		var userReq UserRequest
// 		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
// 			log.Printf("fail: json decode, %v\n", err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		defer r.Body.Close()

// 		// 入力のバリデーション
// 		if userReq.Name == "" || len(userReq.Name) > 50 || userReq.Age < 20 || userReq.Age > 80 {
// 			log.Println("fail: validation error")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		// トランザクションの開始
// 		tx, err := db.Begin()
// 		if err != nil {
// 			log.Printf("fail: db.Begin, %v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// 新しいユーザーIDをULIDで生成
// 		id := ulid.Make().String()
// 		_, err = tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", id, userReq.Name, userReq.Age)
// 		if err != nil {
// 			log.Printf("fail: db.Exec, %v\n", err)
// 			if err := tx.Rollback(); err != nil {
// 				log.Printf("fail: tx.Rollback, %v\n", err)
// 			}
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		if err := tx.Commit(); err != nil {
// 			log.Printf("fail: tx.Commit, %v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// 新しいユーザーIDを返す
// 		response := map[string]string{"id": id}
// 		bytes, err := json.Marshal(response)
// 		if err != nil {
// 			log.Printf("fail: json.Marshal, %v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bytes)

// 	default:
// 		log.Printf("fail: HTTP Method is %s\n", r.Method)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// }

// func main() {
// 	// ② /userでリクエストされたらnameパラメーターと一致する名前を持つレコードをJSON形式で返す
// 	http.HandleFunc("/user", handler)

// 	// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
// 	closeDBWithSysCall()

// 	// 8000番ポートでリクエストを待ち受ける
// 	log.Println("Listening...")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
// func closeDBWithSysCall() {
// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
// 	go func() {
// 		s := <-sig
// 		log.Printf("received syscall, %v", s)

// 		if err := db.Close(); err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Printf("success: db.Close()")
// 		os.Exit(0)
// 	}()
// }
package main

import(
	"net/http"
	"myproject/controllers" // プロジェクト名に応じてパスを変更
	"myproject/dao"
	"myproject/usecase"
	"log"
	"syscall"
	"os"
	"database/sql"
	"os/signal"
)

func main(){
	// DB接続の設定
	db, err := dao.NewDBConnection()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	closeDBWithSysCall(db)

	tweetDao := dao.NewTweetDAO(db)
	// UseCaseのインスタンスを作成
	registerUserUseCase := usecase.NewRegisterUserUseCase(tweetDao) // 適切なDAOを渡す

	// コントローラーのインスタンスを作成
	registerUserController := controllers.NewRegisterUserController(registerUserUseCase)

	// ルート設定
	controllers.RootingRegister(registerUserController)

	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
	
}

// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
func closeDBWithSysCall(db *sql.DB) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}