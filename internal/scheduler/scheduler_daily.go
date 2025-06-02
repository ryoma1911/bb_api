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

// 日次スケジューラをここで設定
func StartDailyFetch(c *cron.Cron) (cron.EntryID, error) {
	id, err := c.AddFunc("30 0 * * *", func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("panic recovered in cron task:", r)
			}
		}()
		log.Println("Executing task at:", time.Now())
		err := GetMatchScheduletoday()
		if err != nil {
			log.Println("Failed task at:", time.Now(), err)
		}
		log.Println("Next task GetMatchScheduletoday:", c.Entries())
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 当日の試合情報を取得しテーブルに登録
func GetMatchScheduletoday() error {
	todate := time.Now().Format("2006-01-02")
	url := "https://baseball.yahoo.co.jp/npb/schedule/?date=" + todate

	res, err := scraper.GetURL(url)
	if err != nil {
		log.Println(fmt.Errorf("failed to get URL: %w", err))
		return err
	}

	doc, err := scraper.GetBody(res)
	if err != nil {
		log.Println(fmt.Errorf("failed to get body: %w", err))
		return err
	}

	matches, err := fetcher.GetMatchSchedule(doc)
	if err != nil {
		log.Println(fmt.Errorf("failed to get match schedule: %w", err))
		return err
	}

	// 試合がある場合はテーブルに格納
	if len(matches) != 0 {
		query_matches := `
			INSERT INTO matches (date, home, away, stadium, starttime, link, league) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
			`
		query_scores := `
			INSERT INTO scores (match_id)
			VALUES (?)
			`
		// DB接続
		db, err := connect.ConnectOnly()
		if err != nil {
			log.Println(fmt.Errorf("failed to check to connect database: %w", err))
			return err
		}
		for _, match := range matches {
			// matchesテーブルに追加
			id, err := repo.InsertData(db, query_matches, match[0], match[1], match[2], match[3], match[5], match[6], match[7])
			if err != nil {
				return err
			}
			// scoresテーブルに追加
			_, err = repo.InsertData(db, query_scores, id)
			if err != nil {
				return err
			}
			log.Println("Get Match Today:", id, match[1], "vs", match[2])
		}
	} else {
		log.Println("There's no game today", time.Now())
	}
	log.Println("Get matches", len(matches), "games")

	return err
}
