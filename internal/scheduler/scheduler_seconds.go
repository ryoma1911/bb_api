package scheduler

import (
	"baseball_report/internal/fetcher"
	"fmt"
	"log"
)

// 試合進捗を取得しテーブル更新
func GetScore(url string, id int) error {
	// DB接続
	db, err := connect.ConnectOnly()
	if err != nil {
		log.Println(fmt.Errorf("failed to check to connect database: %w", err))
		return err
	}
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
	score, err := fetcher.GetMatchScore(doc)
	if err != nil {
		log.Println(fmt.Errorf("failed to get score: %w", err))
		return err
	}
	query := `
			UPDATE scores SET inning = ?, home_score = ?, away_score = ?, batter = ?, result = ? WHERE match_id = ?
		`
	_, err = repo.UpdateData(db, query, score[0], score[1], score[2], score[3], score[4], id)
	if err != nil {
		log.Println(fmt.Errorf("failed to update score: %w", err))
		return err
	}

	return nil
}
