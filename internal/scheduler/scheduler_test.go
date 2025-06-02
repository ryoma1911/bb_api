package scheduler

import (
	"bytes"
	"database/sql"
	"errors"
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
	MockConnectOnly func() (*sql.DB, error)
}

func (m *MockDBHandler) ConnectOnly() (*sql.DB, error) {
	if m.MockConnectOnly != nil {
		return m.MockConnectOnly()
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
	query_match := `
	INSERT INTO matches (date, home, away, stadium, starttime, link, league) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	query_score := `
	INSERT INTO scores (match_id)
	VALUES (?)
	`
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
						<div class="bb-score__link">試合前</div>
						<div class="bb-score__status">12:00</div>
						<div class="bb-score__content" href="test1/index"></div>
					</div>
					<div class="bb-score__item">
						<div class="bb-score__homeLogo">Fighters</div>
						<div class="bb-score__awayLogo">Hawks</div>
						<div class="bb-score__venue">escon</div>
						<div class="bb-score__link">試合前</div>
						<div class="bb-score__status">18:00</div>
						<div class="bb-score__content" href="test2/index"></div>
					</div>
				</div>`))
				return doc, nil
			},
		}

		//DSN取得とDB接続をモック化
		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(query_match).
					WithArgs(todate, "Lions", "Giants", "beruna", "12:00", "test1/score", "Interleague").
					WillReturnResult(sqlmock.NewResult(1, 1)) // match_id=1

				mock.ExpectExec(query_score).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(query_match).
					WithArgs(todate, "Fighters", "Hawks", "escon", "18:00", "test2/score", "Interleague").
					WillReturnResult(sqlmock.NewResult(2, 1)) // match_id=2

				mock.ExpectExec(query_score).
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(2, 1))

				return db, nil
			},
		}

		//関数実行
		err := GetMatchScheduletoday()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Get matches")

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
		err := GetMatchScheduletoday()
		assert.NoError(t, err)

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
			MockConnectOnly: func() (*sql.DB, error) {
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
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec(query).WillReturnError(errors.New("DB insert failed"))
				return db, nil
			},
		}

		GetMatchScheduletoday()

		assert.Contains(t, buf.String(), "failed")
	})
}

func TestGetscore(t *testing.T) {
	// ログ出力のキャプチャ
	var buf bytes.Buffer
	log.SetOutput(&buf)

	t.Run("Success Update Score", func(t *testing.T) {

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
				<body>
					<!-- イニング情報 -->
					<div class="live">
						<em>2回裏</em>
					</div>

					<!-- スコア情報 -->
					<div class="score">
						<table>
							<tr>
								<td class="nm act">オ</td>
								<td>1</td>
							</tr>
							<tr>
								<td class="nm">デ</td>
								<td>2</td>
							</tr>
						</table>
					</div>

					<!-- 打者情報 -->
					<table id="batt">
						<tr>
							<td><a href="/player1">山田</a></td>
						</tr>
					</table>

					<!-- 進捗情報 -->
					<div id="result">
						左2塁打
					</div>
				</body>
				`))
				return doc, nil
			},
		}
		//DSN取得とDB接続をモック化
		// 日付と時間を整形
		todate := time.Now().Format("2006-01-02")
		starttime := time.Now().Format("15:04:05")

		query_match := `
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

		query_score := `
	UPDATE scores SET home_score = ?, away_score = ?, batter = ?, inning = ?, result = ? WHERE match_id = ?
`

		connect = &MockDBHandler{
			MockConnectOnly: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				// SELECT クエリのモック
				mock.ExpectQuery(query_match).
					WithArgs(todate, starttime).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "date", "home", "away", "league", "stadium", "starttime", "link", "inning",
					}).AddRow(
						1, todate, "Lions", "Giants", "Interleague", "beruna", "12:00", "test1/score", "試合前",
					))

				// UPDATE クエリのモック
				mock.ExpectExec(query_score).
					WithArgs("1", "2", "山田", "2回裏", "左2塁打", "1"). // match["id"] は int → 文字列に変換されている
					WillReturnResult(sqlmock.NewResult(1, 1))

				return db, nil
			},
		}

		//対象の関数を実行
		err := GetScores()
		assert.NoError(t, err)

		assert.Contains(t, buf.String(), "Updated Score:")
	})
	t.Run("Error_GetURL", func(t *testing.T) {
		scraper = &MockURLHandler{
			MockGetURL: func(url string) (*http.Response, error) {
				return nil, errors.New("failed to fetch URL")
			},
		}
		//対象の関数を実行
		err := GetScores()
		assert.Error(t, err)

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

		//対象の関数を実行
		err := GetScores()
		assert.Error(t, err)

		assert.Contains(t, buf.String(), "failed to get body: failed to parse HTML")
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
			MockConnectOnly: func() (*sql.DB, error) {
				return nil, errors.New("failed to connect to DB")
			},
		}

		//対象の関数を実行
		err := GetScores()
		assert.Error(t, err)

		assert.Contains(t, buf.String(), "failed to check to connect database: failed to connect to DB")
	})
}

func TestGetscore2(t *testing.T) {
	GetScores()
}