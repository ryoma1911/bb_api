package main

import (
	"baseball_report/internal/api"
	"baseball_report/internal/scheduler"
	"log"
	"net/http"
)

func Run() error {
	//スケジューラ起動
	go scheduler.StartDailyFetch()

	//APIルータを取得しサーバ起動
	router := api.SetupRouter()
	log.Println("API Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	return http.ListenAndServe(":8080", router)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
