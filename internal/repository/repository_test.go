package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestInsertData： DBの特定テーブルにデータが追加されるケースをテスト
func TestInsertData(t *testing.T) {
	repo := &DefaultRepository{}

	// SQLモックの作成
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success INSERT", func(t *testing.T) {
		query := "INSERT INTO matches (team1, team2, score) VALUES (?, ?, ?)"
		mock.ExpectExec(query).
			WithArgs("Yankees", "Red Sox", "5-3").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertData(db, query, "Yankees", "Red Sox", "5-3")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed INSERT", func(t *testing.T) {
		query := "INSERT INTO matches (team1, team2, score) VALUES (?, ?, ?)"
		mock.ExpectExec(query).
			WithArgs("Yankees", "Red Sox", "5-3").
			WillReturnError(sql.ErrConnDone)

		err := repo.InsertData(db, query, "Yankees", "Red Sox", "5-3")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success UPSERT", func(t *testing.T) {
		query := `
			INSERT INTO scores (home_score, away_score, batter, inning, result, match_id)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			home_score = VALUES(home_score),
			away_score = VALUES(away_score),
			batter = VALUES(batter),
			result = VALUES(result)
		`
		mock.ExpectExec(query).
			WithArgs("2", "1", "山田", "3回裏", "ホームラン", 7).
			WillReturnResult(sqlmock.NewResult(1, 2)) // INSERT or UPDATE

		err := repo.InsertData(db, query, "2", "1", "山田", "3回裏", "ホームラン", 7)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed UPSERT", func(t *testing.T) {
		query := `
			INSERT INTO scores (home_score, away_score, batter, inning, result, match_id)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			home_score = VALUES(home_score),
			away_score = VALUES(away_score),
			batter = VALUES(batter),
			result = VALUES(result)
		`
		mock.ExpectExec(query).
			WithArgs("2", "1", "山田", "3回裏", "ホームラン", 7).
			WillReturnError(sql.ErrConnDone)

		err := repo.InsertData(db, query, "2", "1", "山田", "3回裏", "ホームラン", 7)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
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
		query := "SELECT id, date, home, away, league, stadium, starttime, link FROM matches WHERE date ='" + todate + "'"

		//モックの結果を定義
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "link"}).
			AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00", "https://example.com/match/1").
			AddRow(2, todate, "Dodgers", "Giants", "パ・リーグ", "Dodger Stadium", "18:30", "https://example.com/match/2")

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
				"link":      "https://example.com/match/2",
			},
		}
		assert.Equal(t, expected, result)

		// モックの期待値がすべて満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	//試合が始まっているレコードを取得
	t.Run("Success to get match for starttime", func(t *testing.T) {
		//クエリ実行でテーブルからデータが取得されていること
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime, link FROM matches WHERE date ='" + todate + "'" + " AND starttime < CURTIME()"

		//モックの結果を定義
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "link"}).
			AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "17:00", "https://example.com/match/1")

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
				"starttime": "17:00",
				"link":      "https://example.com/match/1",
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

func TestGetMatchAPI(t *testing.T) {
	repo := &DefaultRepository{}

	// SQLモックの作成
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success to get match", func(t *testing.T) {
		//クエリ実行でテーブルからデータが取得されていること
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		//モックの結果を定義
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime"}).
			AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00").
			AddRow(2, todate, "Dodgers", "Giants", "パ・リーグ", "Dodger Stadium", "18:30")

			// モックの期待値を設定
		mock.ExpectQuery(query).WillReturnRows(rows)

		// 関数を実行
		result, err := repo.GetMatchAPI(db, todate)

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
			},
			{
				"id":        2,
				"date":      todate,
				"home":      "Dodgers",
				"away":      "Giants",
				"league":    "パ・リーグ",
				"stadium":   "Dodger Stadium",
				"starttime": "18:30",
			},
		}
		assert.Equal(t, expected, result)

		// モックの期待値がすべて満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	// Failed to get match
	t.Run("Failed to get match", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		// クエリ実行時にエラーを返す
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("query failed"))

		// 関数を実行
		result, err := repo.GetMatchAPI(db, todate)

		// エラーが期待通りであることを確認
		assert.Error(t, err)
		assert.Nil(t, result) // 結果はnilであるべき
		assert.Contains(t, err.Error(), "failed to fetch match")

		// モックの期待値が満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 行のスキャン失敗パターン
	t.Run("Failed to scan", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		// 不正なデータ（型不一致）を返すモック
		rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "stadium", "starttime", "status"}).
			AddRow("invalid_id", time.Now(), "Yankees", "Red Sox", "Yankee Stadium", "19:00", "Scheduled")

		mock.ExpectQuery(query).WillReturnRows(rows)

		// 関数を実行
		result, err := repo.GetMatchAPI(db, todate)

		// エラーが期待通りであることを確認
		assert.Error(t, err)
		assert.Nil(t, result) // 結果はnilであるべき
		assert.Contains(t, err.Error(), "failed to scan match row")

		// モックの期待値が満たされていることを確認
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetScore(t *testing.T) {
	repo := &DefaultRepository{}

	// SQLモックの作成
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success to get score", func(t *testing.T) {
		matchID := "7"
		query := "SELECT id, home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + matchID + "'"

		rows := sqlmock.NewRows([]string{"id", "home_score", "away_score", "batter", "inning", "result", "match_id"}).
			AddRow(1, "2", "1", "山田", "3回裏", "ホームラン", 7)

		mock.ExpectQuery(query).WillReturnRows(rows)

		result, err := repo.GetScore(db, matchID)
		assert.NoError(t, err)

		expected := []map[string]interface{}{
			{
				"id":         1,
				"home_score": "2",
				"away_score": "1",
				"batter":     "山田",
				"inning":     "3回裏",
				"result":     "ホームラン",
				"match_id":   7,
			},
		}

		assert.Equal(t, expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail to get score", func(t *testing.T) {
		matchID := "7"
		query := "SELECT id, home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + matchID + "'"

		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		result, err := repo.GetScore(db, matchID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
