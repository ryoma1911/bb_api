package scheduler

import (
	"baseball_report/internal/fetcher"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

func StartMinutesFetch(c *cron.Cron) (cron.EntryID, error) {
	id, err := c.AddFunc("* 12 * * *", func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("panic recovered in cron task:", r)
			}
		}()
		log.Println("Executing task at:", time.Now())
		err := GetScores()
		if err != nil {
			log.Println("Failed task at:", time.Now(), err)
		}
		log.Println("Next task GetScoreSchedule:", c.Entries())
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 試合進捗を取得しテーブル更新
func GetScores() error {
	// DB接続
	db, err := connect.ConnectOnly()
	if err != nil {
		log.Println(fmt.Errorf("failed to check to connect database: %w", err))
		return err
	}

	//開始中の試合情報を取得
	matches, err := repo.GetMatchScoreLive(db)
	if err != nil {
		log.Println(fmt.Errorf("failed to get to match: %w", err))
		return err
	}
	log.Println("Get Matching :", len(matches))
	for _, match := range matches {
		//試合速報からデータを取得
		res, err := scraper.GetURL(match["link"].(string))
		if err != nil {
			log.Println(fmt.Errorf("failed to get URL: %w", err))
			return err
		}

		doc, err := scraper.GetBody(res)
		if err != nil {
			log.Println(fmt.Errorf("failed to get body: %w", err))
			return err
		}

		score, err := fetcher.GetMatchScore(doc)
		if err != nil {
			log.Println(fmt.Errorf("failed to get match score: %w", err))
			return err
		}
		query := `
				UPDATE scores SET home_score = ?, away_score = ?, batter = ?, inning = ?, result = ? WHERE match_id = ?
				`
		idInt := match["id"].(int)
		idStr := strconv.Itoa(idInt)
		id, err := repo.UpdateData(db, query, score[0][1], score[0][2], score[0][3], score[0][0], score[0][4], idStr)
		if err != nil {
			log.Println(fmt.Errorf("failed to upfate : %w", err))
			return err
		}
		log.Println("Updated Score:", id, score[0][1], "-", score[0][2], score[0][3], score[0][0], score[0][4])
		log.Println("Sleeping 10 second")
		time.Sleep(10 * time.Second)
	}
	return nil
}
