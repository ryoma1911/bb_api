package repository

import (
	"database/sql"
	"fmt"
)

// Repository インターフェース
type Repository interface {
	GetMatch(db *sql.DB, query string) ([]map[string]interface{}, error)
	InsertData(db *sql.DB, query string, args ...interface{}) (int, error)
	UpdateData(db *sql.DB, query string, args ...interface{}) (int, error)
	GetMatchScoreLive(db *sql.DB, todate string, starttime string) ([]map[string]interface{}, error)
}

// DefaultRepository 実装
type DefaultRepository struct{}

func (d *DefaultRepository) InsertData(db *sql.DB, query string, args ...interface{}) (int, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}

func (d *DefaultRepository) UpdateData(db *sql.DB, query string, args ...interface{}) (int, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert: %w", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}

// バックエンド側でDB検索する際に使用
func (d *DefaultRepository) GetMatch(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match: %w", err)
	}
	defer rows.Close()

	var matches []map[string]interface{} //空のスライスを定義
	for rows.Next() {
		var id int
		var date string
		var home string
		var away string
		var league string
		var stadium string
		var starttime string
		var link string
		if err := rows.Scan(&id, &date, &home, &away, &league, &stadium, &starttime, &link); err != nil {
			return nil, fmt.Errorf("failed to scan match row: %w", err)
		}
		//試合情報をマップに格納、スライスに追加
		matches = append(matches, map[string]interface{}{
			"id":        id,
			"date":      date,
			"home":      home,
			"away":      away,
			"league":    league,
			"stadium":   stadium,
			"starttime": starttime,
			"link":      link,
		})

	}
	return matches, nil

}

// 試合情報API出力
func (d *DefaultRepository) GetMatchAPI(db *sql.DB, todate string) ([]map[string]interface{}, error) {
	query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match: %w", err)
	}
	defer rows.Close()

	var matches []map[string]interface{} //空のスライスを定義
	for rows.Next() {
		var id int
		var date string
		var home string
		var away string
		var league string
		var stadium string
		var starttime string
		if err := rows.Scan(&id, &date, &home, &away, &league, &stadium, &starttime); err != nil {
			return nil, fmt.Errorf("failed to scan match row: %w", err)
		}
		//試合情報をマップに格納、スライスに追加
		matches = append(matches, map[string]interface{}{
			"id":        id,
			"date":      date,
			"home":      home,
			"away":      away,
			"league":    league,
			"stadium":   stadium,
			"starttime": starttime,
		})

	}
	return matches, nil

}

// スコア情報を取得
func (d *DefaultRepository) GetScore(db *sql.DB, id string) ([]map[string]interface{}, error) {
	query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + id + "'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match: %w", err)
	}
	defer rows.Close()

	var score []map[string]interface{} //空のスライスを定義
	for rows.Next() {
		var match_id int
		var home_score string
		var away_score string
		var batter string
		var inning string
		var result string

		if err := rows.Scan(&home_score, &away_score, &batter, &inning, &result, &match_id); err != nil {
			return nil, fmt.Errorf("failed to scan match row: %w", err)
		}
		//試合情報をマップに格納、スライスに追加
		score = append(score, map[string]interface{}{
			"match_id":   match_id,
			"home_score": home_score,
			"away_score": away_score,
			"batter":     batter,
			"inning":     inning,
			"result":     result,
		})
	}
	return score, nil

}

// スコア情報を取得
func (d *DefaultRepository) GetMatchScoreLive(db *sql.DB, todate string, starttime string) ([]map[string]interface{}, error) {
	query := `
			SELECT 
				m.id, 
				m.date, 
				m.home, 
				m.away, 
				m.league, 
				m.stadium, 
				m.starttime,
				m.link,
				s.inning
			FROM 
				matches m
			LEFT JOIN 
				scores s ON m.id = s.match_id
			WHERE 
				m.date = ? AND
				m.starttime <= ? AND
				(s.inning <> '試合終了' AND s.inning <> '試合中止')
			`
	rows, err := db.Query(query, todate, starttime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match: %w", err)
	}
	defer rows.Close()

	var matches []map[string]interface{} //空のスライスを定義
	for rows.Next() {
		var id int
		var date string
		var home string
		var away string
		var league string
		var stadium string
		var starttime string
		var link string
		var result string
		if err := rows.Scan(&id, &date, &home, &away, &league, &stadium, &starttime, &link, &result); err != nil {
			return nil, fmt.Errorf("failed to scan match row: %w", err)
		}
		//試合情報をマップに格納、スライスに追加
		matches = append(matches, map[string]interface{}{
			"id":        id,
			"date":      date,
			"home":      home,
			"away":      away,
			"league":    league,
			"stadium":   stadium,
			"starttime": starttime,
			"link":      link,
			"result":    result,
		})

	}
	return matches, nil

}
