package api

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockDBHandler struct {
	MockConnectOnly func() (*sql.DB, error)
}

func (m *MockDBHandler) ConnectOnly() (*sql.DB, error) {
	if m.MockConnectOnly != nil {
		return m.MockConnectOnly()
	}
	db, _, _ := sqlmock.New() // デフォルト動作
	return db, nil
}

// GetMatchesHandler:正常系のパターン
func TestGetMatchesHandler_Success(t *testing.T) {
	// 1リーグ2ゲーム
	t.Run("Get 1league2games", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30")

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		//期待値を設定
		expected := `{
			"セ・リーグ": [
			{
				"id": 1,
				"date": "` + todate + `",
				"home": "Yankees",
				"away": "Red Sox",
				"league": "セ・リーグ",
				"stadium": "Yankee Stadium",
				"starttime": "19:00"
			},
			{
				"id": 2,
				"date": "` + todate + `",
				"home": "Dodgers",
				"away": "Giants",
				"league": "セ・リーグ",
				"stadium": "Dodger Stadium",
				"starttime": "18:30"
			}
			]
		}`

		//HTTPリクエスト作成
		req, _ := http.NewRequest("GET", "/matches", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		// HTTPステータスコードチェック
		assert.Equal(t, http.StatusOK, rr.Code)

		//JSONレスポンスをチェック
		assert.JSONEq(t, expected, rr.Body.String(), "JSON does not match")
	})

	// 2リーグ4ゲーム
	t.Run("Get 2league2games", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30").
					AddRow(3, todate, "SoftBank", "Rakuten", "パ・リーグ", "PayPayドーム", "18:00").
					AddRow(4, todate, "Lotte", "Seibu", "パ・リーグ", "ZOZOマリン", "18:00")

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		//期待値を設定
		// 期待値を設定
		// 期待値を設定
		expected := `{
			"セ・リーグ": [
			{
				"id": 1,
				"date": "` + todate + `",
				"home": "Yankees",
				"away": "Red Sox",
				"league": "セ・リーグ",
				"stadium": "Yankee Stadium",
				"starttime": "19:00"
			},
			{
				"id": 2,
				"date": "` + todate + `",
				"home": "Dodgers",
				"away": "Giants",
				"league": "セ・リーグ",
				"stadium": "Dodger Stadium",
				"starttime": "18:30"
			}
			],
			"パ・リーグ": [
			{
				"id": 3,
				"date": "` + todate + `",
				"home": "SoftBank",
				"away": "Rakuten",
				"league": "パ・リーグ",
				"stadium": "PayPayドーム",
				"starttime": "18:00"
			},
			{
				"id": 4,
				"date": "` + todate + `",
				"home": "Lotte",
				"away": "Seibu",
				"league": "パ・リーグ",
				"stadium": "ZOZOマリン",
				"starttime": "18:00"
			}
			]
		}`

		//HTTPリクエスト作成
		req, _ := http.NewRequest("GET", "/matches", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		// HTTPステータスコードチェック
		assert.Equal(t, http.StatusOK, rr.Code)

		//JSONレスポンスをチェック
		assert.JSONEq(t, expected, rr.Body.String(), "JSON does not match")
	})

	// 1試合もない
	t.Run("Get Nogames", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime"})
				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		//期待値を設定
		expected := `{
			"message": "No matches found"
		}`

		//HTTPリクエスト作成
		req, _ := http.NewRequest("GET", "/matches", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		// HTTPステータスコードチェック
		assert.Equal(t, http.StatusOK, rr.Code)

		//JSONレスポンスをチェック
		assert.JSONEq(t, expected, rr.Body.String(), "JSON does not match")
	})

}

// GetMatchesHandler:エラー系のパターン
func TestGetMatchesHandler_Failes(t *testing.T) {

	// DB接続失敗
	t.Run("Failed to connect DB", func(t *testing.T) {
		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				return nil, errors.New("DB接続エラー")
			},
		}

		req, _ := http.NewRequest("GET", "/matches", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	// クエリ実行失敗
	t.Run("Failed to execute query", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(query).WillReturnError(errors.New("クエリエラー"))
				return db, nil
			},
		}

		req, _ := http.NewRequest("GET", "/matches", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error executing query:")
	})
}

// SetupRouter:正常パターン
func TestSetupRouter_Success(t *testing.T) {
	router := SetupRouter()

	t.Run("GET /matches returns match data", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30")

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		req := httptest.NewRequest("GET", "/matches", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)

		bodyBytes, _ := io.ReadAll(rr.Body)

		assert.Contains(t, string(bodyBytes), "Giants", "response should contain 'Giants'")
	})

	t.Run("GET /scores returns score data", func(t *testing.T) {
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='7'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"home_score", "away_score", "batter", "inning", "result", "match_id"}).
					AddRow("2", "1", "山田", "3回裏", "ホームラン", 7)

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		req := httptest.NewRequest("GET", "/scores/7", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)

		bodyBytes, _ := io.ReadAll(rr.Body)

		assert.Contains(t, string(bodyBytes), "7", "response should contain '7'")

	})
	t.Run("GET /health returns OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rr.Code)
		}

		if rr.Body.String() != "OK" {
			t.Errorf("expected body OK, got %s", rr.Body.String())
		}
	})
}

func TestGetScoreHandler_Success(t *testing.T) {
	// 取得成功
	t.Run("Success get score", func(t *testing.T) {
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='7'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"home_score", "away_score", "batter", "inning", "result", "match_id"}).
					AddRow("2", "1", "山田", "3回裏", "ホームラン", 7)

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		//期待値を設定
		expected := `[
				{
					"home_score": "2",
					"away_score": "1",
					"batter": "山田",
					"inning": "3回裏",
					"result": "ホームラン",
					"match_id": 7
				}
				]`

		//HTTPリクエスト作成
		req := httptest.NewRequest("GET", "/scores/7", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetScoreHandler)

		handler.ServeHTTP(rr, req)

		// HTTPステータスコードチェック
		assert.Equal(t, http.StatusOK, rr.Code)

		//JSONレスポンスをチェック
		assert.JSONEq(t, expected, rr.Body.String(), "JSON does not match")
	})

	t.Run("Success no score", func(t *testing.T) {
		query := "SELECT home_score, away_score, batter, inning, result, match_id FROM scores WHERE match_id ='7'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"home_score", "away_score", "batter", "inning", "result", "match_id"})

				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}
		//期待値を設定
		expected := `{
				"message": "No score found"
			}`

		//HTTPリクエスト作成
		req := httptest.NewRequest("GET", "/scores/7", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetScoreHandler)

		handler.ServeHTTP(rr, req)

		// HTTPステータスコードチェック
		assert.Equal(t, http.StatusOK, rr.Code)

		//JSONレスポンスをチェック
		assert.JSONEq(t, expected, rr.Body.String(), "JSON does not match")
	})
}

// GetMatchesHandler:エラー系のパターン
func TestGetScoreHandler_Failes(t *testing.T) {
	// DB接続失敗
	t.Run("Failed to connect DB", func(t *testing.T) {
		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				return nil, errors.New("DB接続エラー")
			},
		}

		req, _ := http.NewRequest("GET", "/scores/7", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetScoreHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	// クエリ実行失敗
	t.Run("Failed to execute query", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(query).WillReturnError(errors.New("クエリエラー"))
				return db, nil
			},
		}

		req, _ := http.NewRequest("GET", "/scores/7", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetScoreHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error executing query:")
	})
}
