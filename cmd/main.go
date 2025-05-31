package main

import (
	"baseball_report/internal/api"
	"baseball_report/internal/scheduler"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

func Run() error {

	//スケジューラ起動
	location, _ := time.LoadLocation("Asia/Tokyo")
	c := cron.New(cron.WithLocation(location))

	// スケジューラ登録
	go func() {
		id, err := scheduler.StartDailyFetch(c)
		if err != nil {
			log.Println("Failed to register cron job:", err)
			return
		}
		// スケジューラ開始
		c.Start()

		// 次回実行時刻を取得してログに出力
		entry := c.Entry(id)
		log.Printf("Cron job registered! Next scheduled run: %s (JST)", entry.Next)
		select {}
	}()

	//APIルータを取得しサーバ起動
	router := api.SetupRouter()
	log.Println("API Server running")
	log.Println("Now time is: ", time.Now())

	return http.ListenAndServe(":8080", router)
}

func main() {
	if err := Run(); err != nil {

		log.Fatal(err)
	}
}
