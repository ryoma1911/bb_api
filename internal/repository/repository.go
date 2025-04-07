package repository

import (
	"database/sql"
	"fmt"
)

// Repository インターフェース
type Repository interface {
	GetMatch(db *sql.DB, query string) ([]map[string]interface{}, error)
	InsertMatch(db *sql.DB, query string, args ...interface{}) error
}

// DefaultRepository 実装
type DefaultRepository struct{}

func (d *DefaultRepository) InsertMatch(db *sql.DB, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
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
