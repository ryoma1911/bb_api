package scheduler

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
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

type MockURLHandler struct {
	MockGetURL  func(url string) (*http.Response, error)
	MockGetBody func(res *http.Response) (*goquery.Document, error)
}

func (m *MockURLHandler) GetURL(url string) (*http.Response, error) {
	if m.MockGetURL != nil {
		return m.MockGetURL(url)
	}
	return nil, nil
}

func (m *MockURLHandler) GetBody(res *http.Response) (*goquery.Document, error) {
	if m.MockGetBody != nil {
		return m.MockGetBody(res)
	}
	return nil, nil
}

func TestStartDailyFetch_Success(t *testing.T) {
	// ログ出力のキャプチャ
	var buf bytes.Buffer
	log.SetOutput(&buf)
	c := cron.New(cron.WithLocation(time.Local))

	id, err := StartDailyFetch(c)
	assert.NoError(t, err)

	c.Start()
	defer c.Stop()

	time.Sleep(500 * time.Millisecond)

	entry := c.Entry(id)
	assert.Equal(t, entry.ID, id)

	assert.False(t, entry.Next.IsZero())
}

// 正常系のパターン
func TestGetMatchScheduletoday_Success(t *testing.T) {
	todate := time.Now().Format("2006/01/02")
	query := "INSERT INTO matches (date, home, away, stadium, starttime, link, league) VALUES (?, ?, ?, ?, ?, ?, ?)"
	// ログ出力のキャプチャ
	var buf bytes.Buffer
	log.SetOutput(&buf)

	//1リーグ2ゲームの場合
	t.Run("Get 1league2games", func(t *testing.T) {

		//スクレイピング処理をモック化
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("mock html")),
				}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
				<div class="bb-score">
					<h2 class="bb-score__title">Interleague</h2>
					<div class="bb-score__item">
						<div class="bb-score__homeLogo">Lions</div>
						<div class="bb-score__awayLogo">Giants</div>
						<div class="bb-score__venue">beruna</div>
						<div class="bb-score__status">試合前</div>
						<div class="bb-score__link">12:00</div>
						<div class="bb-score__content" href="test1/index"></div>
					</div>
					<div class="bb-score__item">
						<div class="bb-score__homeLogo">Fighters</div>
						<div class="bb-score__awayLogo">Hawks</div>
						<div class="bb-score__venue">escon</div>
						<div class="bb-score__status">試合前</div>
						<div class="bb-score__link">18:00</div>
						<div class="bb-score__content" href="test2/index"></div>
					</div>
				</div>`))
				return doc, nil
			},
		}

		//DSN取得とDB接続をモック化
		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", nil
			},
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(query).
					WithArgs(todate, "Lions", "Giants", "beruna", "12:00", "test1/score", "Interleague").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(query).
					WithArgs(todate, "Fighters", "Hawks", "escon", "18:00", "test2/score", "Interleague").
					WillReturnResult(sqlmock.NewResult(2, 2))

				return db, nil
			},
		}

		//関数実行
		GetMatchScheduletoday()
		assert.Contains(t, buf.String(), "Get matches 2 games")

	})
	t.Run("Get Nogame", func(t *testing.T) {
		//スクレイピング処理をモック化
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("mock html")),
				}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
				<div class="bb-score">
					<div class="bb-noData">今日は試合がありません。</div>
				</div>`))
				return doc, nil
			},
		}
		//関数実行
		GetMatchScheduletoday()

		//ログ結果が期待値と一致している
		assert.Contains(t, buf.String(), "There's no game today")
	})

}

// エラー系のテスト
func TestGetMatchScheduletoday_Errors(t *testing.T) {
	query := "INSERT INTO matches (date, home, away, stadium, status, starttime, link, league) values (?, ?, ?, ?, ?, ?, ?, ?)"

	var buf bytes.Buffer
	log.SetOutput(&buf)

	t.Run("Error_GetURL", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return nil, errors.New("failed to fetch URL")
			},
		}

		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed to get URL: failed to fetch URL")
	})

	t.Run("Error_GetBody", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock html"))}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				return nil, errors.New("failed to parse HTML")
			},
		}

		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed to get body: failed to parse HTML")
	})

	t.Run("Error_GetDSN", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock html"))}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
					<div class="bb-score">
						<div class="bb-score__item">
							<div class="bb-score__homeLogo">Lions</div>
							<div class="bb-score__awayLogo">Giants</div>
							<div class="bb-score__venue">beruna</div>
							<div class="bb-score__status">試合前</div>
							<div class="bb-score__link">12:00</div>
							<div class="bb-score__content" href="test1/index"></div>
						</div>
					</div>`))
				return doc, nil
			},
		}
		connect = &MockDBHandler{
			MockGetDSNFromEnv: func(path string) (string, error) {
				return "mock dsn", fmt.Errorf("failed to load env file:")
			},
		}
		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed to load env file: ")
	})

	t.Run("Error_DBConnect", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock html"))}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
					<div class="bb-score">
						<div class="bb-score__item">
							<div class="bb-score__homeLogo">Lions</div>
							<div class="bb-score__awayLogo">Giants</div>
							<div class="bb-score__venue">beruna</div>
							<div class="bb-score__status">試合前</div>
							<div class="bb-score__link">12:00</div>
							<div class="bb-score__content" href="test1/index"></div>
						</div>
					</div>`))
				return doc, nil
			},
		}

		connect = &MockDBHandler{
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				return nil, errors.New("failed to connect to DB")
			},
		}

		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed to check to connect database: failed to connect to DB")
	})

	t.Run("Error_DBInsert", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock html"))}, nil
			},
			MockGetBody: func(res *http.Response) (*goquery.Document, error) {
				doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
					<div class="bb-score">
						<div class="bb-score__item">
							<div class="bb-score__homeLogo">Lions</div>
							<div class="bb-score__awayLogo">Giants</div>
							<div class="bb-score__venue">beruna</div>
							<div class="bb-score__status">試合前</div>
							<div class="bb-score__link">12:00</div>
							<div class="bb-score__content" href="test1/index"></div>
						</div>
					</div>`))
				return doc, nil
			},
		}

		connect = &MockDBHandler{
			MockConnectOnly: func(dsn string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec(query).WillReturnError(errors.New("DB insert failed"))
				return db, nil
			},
		}

		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed to insert:")
	})
}
