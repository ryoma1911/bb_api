package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestInsertMatch： DBの特定テーブルにデータが追加されるケースをテスト
func TestInsertMatch(t *testing.T) {
	repo := &DefaultRepository{}

	// SQLモックの作成
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success INSERT", func(t *testing.T) {
		//クエリを実行しデータが追加されていること
		query := "INSERT INTO matches (team1, team2, score) VALUES (?, ?, ?)"
		mock.ExpectExec(query).
			WithArgs("Yankees", "Red Sox", "5-3").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.InsertMatch(db, query, "Yankees", "Red Sox", "5-3")
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Failed INSERT", func(t *testing.T) {
		//クエリ実行でエラーが出力されていること
		query := "INSERTE INTO matches (team1, team2, score) VALUES (?, ?)"
		mock.ExpectExec(query).
			WithArgs("Yankees", "Red Sox", "5-3").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.InsertMatch(db, query, "")
		assert.Error(t, err)

		assert.Error(t, mock.ExpectationsWereMet())
	})
}

func TestGetMatch(t *testing.T) {
	repo := &DefaultRepository{}

	// SQLモックの作成
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success to get match", func(t *testing.T) {
		//クエリ実行でテーブルからデータが取得されていること
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches"

		//モックの結果を定義
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "status", "link"}).
			AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00", "Scheduled", "https://example.com/match/1").
			AddRow(2, todate, "Dodgers", "Giants", "パ・リーグ", "Dodger Stadium", "18:30", "Scheduled", "https://example.com/match/2")

			// モックの期待値を設定
		mock.ExpectQuery(query).WillReturnRows(rows)

		// 関数を実行
		result, err := repo.GetMatch(db, query)

		// エラーが発生しないことを確認
		assert.NoError(t, err)
		// 返却結果が期待通りであることを確認
		expected := []map[string]interface{}{
			{
				"id":        1,
				"date":      todate,
				"home":      "Yankees",
				"away":      "Red Sox",
				"league":    "セ・リーグ",
				"stadium":   "Yankee Stadium",
				"starttime": "19:00",
				"status":    "Scheduled",
				"link":      "https://example.com/match/1",
			},
			{
				"id":        2,
				"date":      todate,
				"home":      "Dodgers",
				"away":      "Giants",
				"league":    "パ・リーグ",
				"stadium":   "Dodger Stadium",
				"starttime": "18:30",
				"status":    "Scheduled",
				"link":      "https://example.com/match/2",
			},
		}
		assert.Equal(t, expected, result)

		// モックの期待値がすべて満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	// Failed to get match
	t.Run("Failed to get match", func(t *testing.T) {
		query := "SELECT id, date, home, away, stadium, starttime, status, link FROM matches"

		// クエリ実行時にエラーを返す
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("query failed"))

		// 関数を実行
		result, err := repo.GetMatch(db, query)

		// エラーが期待通りであることを確認
		assert.Error(t, err)
		assert.Nil(t, result) // 結果はnilであるべき
		assert.Contains(t, err.Error(), "failed to fetch match")

		// モックの期待値が満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 行のスキャン失敗パターン
	t.Run("Failed to scan", func(t *testing.T) {
		query := "SELECT id, date, home, away, stadium, starttime, status, link FROM matches"

		// 不正なデータ（型不一致）を返すモック
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "stadium", "starttime", "status", "link"}).
			AddRow("invalid_id", time.Now(), "Yankees", "Red Sox", "Yankee Stadium", "19:00", "Scheduled", "https://example.com/match/1")

		mock.ExpectQuery(query).WillReturnRows(rows)

		// 関数を実行
		result, err := repo.GetMatch(db, query)

		// エラーが期待通りであることを確認
		assert.Error(t, err)
		assert.Nil(t, result) // 結果はnilであるべき
		assert.Contains(t, err.Error(), "failed to scan match row")

		// モックの期待値が満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
