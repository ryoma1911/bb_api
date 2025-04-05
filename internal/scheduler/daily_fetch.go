package scheduler

import (
	db "baseball_report/internal/config"
	"baseball_report/internal/fetcher"
	"baseball_report/internal/repository"
	"baseball_report/utils"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// 依存関係を抽象化するためインターフェース化
var repo repository.Repository = &repository.DefaultRepository{}
var connect db.DBHandler = &db.DBService{}
var scraper utils.URLHandler = &utils.URLService{}

var scheduler *cron.Cron

func StartDailyFetch() {

	scheduler = cron.New()
	url := "https://baseball.yahoo.co.jp/npb/schedule/"

	//毎日6:00に実行
	scheduler.AddFunc("0 6 * * *", func() {
		log.Println("Executing task at:", time.Now())

		res, err := scraper.GetURL(url)
		if err != nil {
			log.Println(fmt.Errorf("failed to get URL: %w", err))
			return
		}

		doc, err := scraper.GetBody(res)
		if err != nil {
			log.Println(fmt.Errorf("failed to get body: %w", err))
			return
		}

		matches, err := fetcher.GetMatchSchedule(doc)
		if err != nil {
			log.Println(fmt.Errorf("failed to get match schedule: %w", err))
			return
		}

		//試合がある場合はテーブルに格納
		if len(matches) != 0 {
			query := "INSERT INTO matches (date, home, away, stadium, status, starttime, link, league) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
			dsn, err := connect.GetDSNFromEnv("/code/.env")
			if err != nil {
				log.Println(fmt.Errorf("failed to load env file: %w", err))
				return
			}
			db, err := connect.ConnectOnly(dsn)
			if err != nil {
				log.Println(fmt.Errorf("failed to check to connect database: %w", err))
				return
			}
			for _, match := range matches {
				err := repo.InsertMatch(db, query, match[0], match[1], match[2], match[3], match[4], match[5], match[6], match[7])
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			log.Println("There's no game today", time.Now())
		}
		log.Println("Get matches", len(matches), "games")

	})

}
