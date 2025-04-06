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
	"github.com/stretchr/testify/assert"
)

type MockDBHandler struct {
	MockGetDSNFromEnv func(path string) (string, error)
	MockConnectOnly   func(dsn string) (*sql.DB, error)
}

func (m *MockDBHandler) GetDSNFromEnv(path string) (string, error) {
	if m.MockGetDSNFromEnv != nil {
		return m.MockGetDSNFromEnv(path)
	}
	return "mock_user:mock_password@tcp(mock_db:3306)/testdb", nil
}

func (m *MockDBHandler) ConnectOnly(dsn string) (*sql.DB, error) {
	if m.MockConnectOnly != nil {
		return m.MockConnectOnly(dsn)
	}
	db, _, _ := sqlmock.New() // デフォルト動作
	return db, nil
}

// GetMatchesHandler:正常系のパターン
func TestGetMatchesHandler_Success(t *testing.T) {
	// 1リーグ2ゲーム
	t.Run("Get 1league2games", func(t *testing.T) {
		todate := time.Now().Format("2006/01/02")
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", nil
			},
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "status", "link"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00", "試合前", "https://example.com/match/1").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30", "試合前", "https://example.com/match/2")

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
				"starttime": "19:00",
				"status": "試合前",
				"link": "https://example.com/match/1"
			  },
			  {
				"id": 2,
				"date": "` + todate + `",
				"home": "Dodgers",
				"away": "Giants",
				"league": "セ・リーグ",
				"stadium": "Dodger Stadium",
				"starttime": "18:30",
				"status": "試合前",
				"link": "https://example.com/match/2"
			  }
			]
		}`

		//HTTPリクエスト作成
		req, err := http.NewRequest("GET", "/matches", nil)
		assert.NoError(t, err)

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
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", nil
			},
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "status", "link"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00", "試合前", "https://example.com/match/1").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30", "試合前", "https://example.com/match/2").
					AddRow(3, todate, "SoftBank", "Rakuten", "パ・リーグ", "PayPayドーム", "18:00", "試合前", "https://example.com/match/3").
					AddRow(4, todate, "Lotte", "Seibu", "パ・リーグ", "ZOZOマリン", "18:00", "試合前", "https://example.com/match/4")

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
				"starttime": "19:00",
				"status": "試合前",
				"link": "https://example.com/match/1"
			},
			{
				"id": 2,
				"date": "` + todate + `",
				"home": "Dodgers",
				"away": "Giants",
				"league": "セ・リーグ",
				"stadium": "Dodger Stadium",
				"starttime": "18:30",
				"status": "試合前",
				"link": "https://example.com/match/2"
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
				"starttime": "18:00",
				"status": "試合前",
				"link": "https://example.com/match/3"
			},
			{
				"id": 4,
				"date": "` + todate + `",
				"home": "Lotte",
				"away": "Seibu",
				"league": "パ・リーグ",
				"stadium": "ZOZOマリン",
				"starttime": "18:00",
				"status": "試合前",
				"link": "https://example.com/match/4"
			}
			]
		}`

		//HTTPリクエスト作成
		req, err := http.NewRequest("GET", "/matches", nil)
		assert.NoError(t, err)

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
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", nil
			},
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "status", "link"})
				mock.ExpectQuery(query).WillReturnRows(rows)
				return db, nil
			},
		}

		//期待値を設定
		expected := `{
			"message": "No matches found"
		}`

		//HTTPリクエスト作成
		req, err := http.NewRequest("GET", "/matches", nil)
		assert.NoError(t, err)

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
	// DSN取得失敗
	t.Run("Failed to get DSN", func(t *testing.T) {
		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(_ string) (string, error) {
				return "", errors.New("DSN取得失敗")
			},
			MockConnectOnly: nil,
		}

		req, _ := http.NewRequest("GET", "/matches", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetMatchesHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Get dsn error")
	})

	// DB接続失敗
	t.Run("Failed to connect DB", func(t *testing.T) {
		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(_ string) (string, error) {
				return "mock_dsn", nil
			},
			MockConnectOnly: func(_ string) (*sql.DB, error) {
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
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(_ string) (string, error) {
				return "mock_dsn", nil
			},
			MockConnectOnly: func(_ string) (*sql.DB, error) {
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
		query := "SELECT id, date, home, away, league, stadium, starttime, status, link FROM matches WHERE date ='" + todate + "'"

		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", nil
			},
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"id", "date", "home", "away", "league", "stadium", "starttime", "status", "link"}).
					AddRow(1, todate, "Yankees", "Red Sox", "セ・リーグ", "Yankee Stadium", "19:00", "試合前", "https://example.com/match/1").
					AddRow(2, todate, "Dodgers", "Giants", "セ・リーグ", "Dodger Stadium", "18:30", "試合前", "https://example.com/match/2")

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
