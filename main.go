// package main

// import(
// 	"net/http"
// 	"myproject/controllers" 
// 	"myproject/middleware"
// 	"github.com/gorilla/mux"
// 	"myproject/dao"
// 	"log"
// 	"syscall"
// 	"os"
// 	"context"
// 	"database/sql"
// 	"os/signal"
// )

// func main(){
// 	// DB接続の設定
// 	db, err := dao.NewDBConnection()
// 	if err != nil {
// 		log.Fatal("Failed to connect to DB: ", err)
// 	}

// 	closeDBWithSysCall(db)

// 	tweetDao := dao.NewTweetDAO(db)
// 	postDAO := dao.NewPostDAO(db)
// 	// コントローラーのインスタンスを作成
// 	userController:= controllers.NewUserController(tweetDao)
// 	postController := controllers.NewPostController(postDAO)

// 	// ルート設定
// 	r := mux.NewRouter()
// 	r.Use(middleware.CORSMiddleware)
// 	userController.RegisterRoutes(r)
// 	postController.RegisterRoutes(r)

// 	log.Println("Listening...")
// 	if err := http.ListenAndServe(":8000", r); err != nil {
// 		log.Fatal(err)
// 	}

// }

// //Ctrl+CでHTTPサーバー停止時にDBをクローズする
// func closeDBWithSysCall(db *sql.DB) {
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

import (
	"context"
	"net/http"
	"database/sql"
	"log"
	"myproject/controllers"
	"myproject/dao"
	"myproject/middleware"
	"os"
	"os/signal"
	"syscall"
	"github.com/gorilla/mux"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	// DB接続の設定
	db, err := dao.NewDBConnection()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	// DBクローズを設定
	closeDBWithSysCall(db)

	// TweetDAOとPostDAOの初期化
	tweetDao := dao.NewTweetDAO(db)
	postDAO := dao.NewPostDAO(db)

	// Gemini DAOの初期化
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal("Failed to create Gemini client: ", err)
	}
	geminiDao := dao.NewVertexAiDAO(client)

	// コントローラーのインスタンスを作成
	userController := controllers.NewUserController(tweetDao)
	postController := controllers.NewPostController(postDAO)
	geminiController := controllers.NewGeminiController(geminiDao)

	// ルート設定
	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)
	userController.RegisterRoutes(r)
	postController.RegisterRoutes(r)
	geminiController.RegisterRoutes(r) // GeminiControllerのルートを追加

	// HTTPサーバーを開始
	log.Println("Listening on :8000...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// Ctrl+CでHTTPサーバー停止時にDBをクローズする
func closeDBWithSysCall(db *sql.DB) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("Received signal: %v", s)

		// DBのクローズ処理
		if err := db.Close(); err != nil {
			log.Fatal("Error closing DB connection: ", err)
		}
		log.Println("DB connection closed successfully")
		os.Exit(0)
	}()
}
