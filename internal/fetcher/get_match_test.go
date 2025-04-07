package fetcher

import (
	"baseball_report/utils"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// Testgetmatchschedulemock： モックから試合情報取得するケースをテスト
func TestGetmatchschedulemock(t *testing.T) {
	todate := time.Now().Format("2006/01/02")

	html1league1games := `
	<div class="bb-score">
		<h2 class="bb-score__title">Interleague</h2>
		<div class="bb-score__item">
			<div class="bb-score__homeLogo">Lions</div>
			<div class="bb-score__awayLogo">Giants</div>
			<div class="bb-score__venue">beruna</div>
			<div class="bb-score__status">試合前</div>
			<div class="bb-score__link">12:00</div>
			<div class="bb-score__content" href="/index"></div>
		</div>
	</div>`

	html1league2games := `
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
	</div>`

	html2league2games := `
	<div class="bb-score">
		<h2 class="bb-score__title">Aleague</h2>
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
	</div>
	<div class="bb-score">
		<h2 class="bb-score__title">Bleague</h2>
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
	</div>`

	htmlnogame := `
	<div class="bb-score">
		<div class="bb-noData">今日は試合がありません。</div>
	</div>`

	// 正常系
	// リーグ1件,試合1件
	t.Run("Success get 1league1game", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{todate, "Lions", "Giants", "beruna", "12:00", "試合前", "/score", "Interleague"},
		}

		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html1league1games))
		matchdata, err := GetMatchSchedule(doc)

		//エラーが発生せず処理が終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Equal(t, expected, matchdata)
	})

	// リーグ1件,試合2件
	t.Run("Success get 1league2game", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{todate, "Lions", "Giants", "beruna", "12:00", "試合前", "test1/score", "Interleague"},
			{todate, "Fighters", "Hawks", "escon", "18:00", "試合前", "test2/score", "Interleague"},
		}

		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html1league2games))
		matchdata, err := GetMatchSchedule(doc)

		//エラーが発生せず処理が終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Equal(t, expected, matchdata)

	})
	// リーグ2件,試合2件
	t.Run("Success get 2league2game", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{todate, "Lions", "Giants", "beruna", "12:00", "試合前", "test1/score", "Aleague"},
			{todate, "Fighters", "Hawks", "escon", "18:00", "試合前", "test2/score", "Aleague"},
			{todate, "Lions", "Giants", "beruna", "12:00", "試合前", "test1/score", "Bleague"},
			{todate, "Fighters", "Hawks", "escon", "18:00", "試合前", "test2/score", "Bleague"},
		}

		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html2league2games))
		matchdata, err := GetMatchSchedule(doc)

		//エラーが発生せず処理が終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Equal(t, expected, matchdata)

	})
	// 今日は試合が無い
	t.Run("Success get nogame", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlnogame))
		matchdata, err := GetMatchSchedule(doc)

		//エラーが発生せず処理が終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Nil(t, matchdata)

	})

}

// Testgetmatchscheduleprod： 本番から試合情報取得するケースをテスト
func TestGetmatchscheduleprod(t *testing.T) {
	scraper := utils.URLService{}
	t.Run("Success get nogame", func(t *testing.T) {
		//取得元のURLを定義
		url := "https://baseball.yahoo.co.jp/npb/schedule/?date=2025-03-27"
		// 野球速報サイトからデータを取得
		res, err := scraper.GetURL(url)
		assert.NoError(t, err)

		doc, err := scraper.GetBody(res)
		assert.NoError(t, err)

		matchdata, err := GetMatchSchedule(doc)
		assert.NoError(t, err)

		assert.Nil(t, matchdata)
	})

}
