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

func StartDailyFetch(c *cron.Cron) (cron.EntryID, error) {
	// 日本時間でcronスケジュールを設定
	return c.AddFunc("10 0 * * *", func() {
		// 現在の日本時間をログに出力
		log.Println("Executing task at:", time.Now())
		GetMatchScheduletoday()
	})
}

// 当日の試合情報を取得しテーブルに登録
func GetMatchScheduletoday() {
	todate := time.Now().Format("2006-01-02")
	url := "https://baseball.yahoo.co.jp/npb/schedule/?date=" + todate

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
		query := `
			INSERT INTO matches (date, home, away, stadium, starttime, link, league) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
			`
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
			//matchesテーブルに追加
			err := repo.InsertData(db, query, match[0], match[1], match[2], match[3], match[5], match[6], match[7])
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		log.Println("There's no game today", time.Now())
	}
	log.Println("Get matches", len(matches), "games")
}
