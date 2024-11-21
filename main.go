package main

import(
	"net/http"
	"myproject/controllers" 
	"github.com/gorilla/mux"
	"myproject/dao"
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
	postDAO := dao.NewPostDAO(db)
	// コントローラーのインスタンスを作成
	userController:= controllers.NewUserController(tweetDao)
	postController := controllers.NewPostController(postDAO)

	// ルート設定
	r := mux.NewRouter()
	userController.RegisterRoutes(r)
	postController.RegisterRoutes(r)

	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}

}

//Ctrl+CでHTTPサーバー停止時にDBをクローズする
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