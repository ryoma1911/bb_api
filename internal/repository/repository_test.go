package repository

import (
	"database/sql"
	"fmt"
	"regexp"
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

		id, err := repo.InsertData(db, query, "Yankees", "Red Sox", "5-3")
		assert.NoError(t, err)
		assert.Equal(t, id, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed INSERT", func(t *testing.T) {
		query := "INSERT INTO matches (team1, team2, score) VALUES (?, ?, ?)"
		mock.ExpectExec(query).
			WithArgs("Yankees", "Red Sox", "5-3").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.InsertData(db, query, "Yankees", "Red Sox", "5-3")
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

		_, err := repo.InsertData(db, query, "2", "1", "山田", "3回裏", "ホームラン", 7)
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

		_, err := repo.InsertData(db, query, "2", "1", "山田", "3回裏", "ホームラン", 7)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success UPDATE", func(t *testing.T) {
		query := `
			UPDATE scores SET home_score = ?, away_score = ?, batter = ?, inning = ?, result = ? WHERE match_id = ?
		`
		mock.ExpectExec(query).
			WithArgs("2", "2", "山田", "3回裏", "ホームラン", 7).
			WillReturnResult(sqlmock.NewResult(1, 2)) // INSERT or UPDATE

		_, err := repo.UpdateData(db, query, "2", "2", "山田", "3回裏", "ホームラン", 7)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed UPSERT", func(t *testing.T) {
		query := `
			UPDATE scores SET home_score = ?, away_score = ?, batter = ?, inning = ?, result = ? WHERE match_id = ?
		`
		mock.ExpectExec(query).
			WithArgs("2", "1", "山田", "3回裏", "ホームラン", 7).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.UpdateData(db, query, "2", "1", "山田", "3回裏", "ホームラン", 7)
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

	t.Run("Success to get score result=試合中", func(t *testing.T) {
		matchID := "7"
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + matchID + "'"

		rows := sqlmock.NewRows([]string{"home_score", "away_score", "batter", "inning", "result", "match_id"}).
			AddRow("2", "1", "山田", "3回裏", "ホームラン", 7)

		mock.ExpectQuery(query).WillReturnRows(rows)

		result, err := repo.GetScore(db, matchID)
		assert.NoError(t, err)

		expected := []map[string]interface{}{
			{
				"match_id":   7,
				"home_score": "2",
				"away_score": "1",
				"batter":     "山田",
				"inning":     "3回裏",
				"result":     "ホームラン",
			},
		}

		assert.Equal(t, expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success to get score result=試合前", func(t *testing.T) {
		matchID := "7"
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + matchID + "'"

		rows := sqlmock.NewRows([]string{"home_score", "away_score", "batter", "inning", "result", "match_id"}).
			AddRow("0", "0", "テスト", "試合前", "試合前", 7)

		mock.ExpectQuery(query).WillReturnRows(rows)

		result, err := repo.GetScore(db, matchID)
		assert.NoError(t, err)

		expected := []map[string]interface{}{
			{
				"match_id":   7,
				"home_score": "0",
				"away_score": "0",
				"batter":     "テスト",
				"inning":     "試合前",
				"result":     "試合前",
			},
		}

		assert.Equal(t, expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail to get score", func(t *testing.T) {
		matchID := "7"
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='" + matchID + "'"

		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		result, err := repo.GetScore(db, matchID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetMatchScoreLive(t *testing.T) {
	repo := &DefaultRepository{}

	// sqlmockの準備
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 日付・開始時刻
	today := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	startTime := time.Date(2025, 6, 1, 13, 5, 0, 0, time.UTC).Format("15:04:05")

	t.Run("Success to get ongoing matches", func(t *testing.T) {
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
		rows := sqlmock.NewRows([]string{
			"id", "date", "home", "away", "league", "stadium", "starttime", "link", "result",
		}).AddRow(1, "2025-06-01", "チームA", "チームB", "セリーグ", "東京ドーム", "13:00:00", "http://example.com", "試合中")

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(today, startTime).WillReturnRows(rows)

		results, err := repo.GetMatchScoreLive(db, today, startTime)
		assert.NoError(t, err)

		expected := []map[string]interface{}{
			{
				"id":        1,
				"date":      "2025-06-01",
				"home":      "チームA",
				"away":      "チームB",
				"league":    "セリーグ",
				"stadium":   "東京ドーム",
				"starttime": "13:00:00",
				"link":      "http://example.com",
				"result":    "試合中",
			},
		}
		assert.Equal(t, expected, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success with no results", func(t *testing.T) {
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

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(today, startTime).WillReturnRows(sqlmock.NewRows([]string{
			"id", "date", "home", "away", "league", "stadium", "starttime", "link", "result",
		}))

		results, err := repo.GetMatchScoreLive(db, today, startTime)
		assert.NoError(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail to query", func(t *testing.T) {
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

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(today, startTime).WillReturnError(sql.ErrConnDone)

		results, err := repo.GetMatchScoreLive(db, today, startTime)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
