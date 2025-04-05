package fetcher

import (
	"baseball_report/utils"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestGetmatchscoremock(t *testing.T) {
	inninghtml := `<body>

    <!-- イニング情報 -->
    <div class="live">
        <em>5回裏</em>
    </div>

    <!-- スコア情報 -->
    <div class="score">
        <table>
            <tr>
                <td class="nm act">オ</td>
                <td>0</td>
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
            <td><a href="/player1">山田 太郎</a></td>
        </tr>
    </table>

    <!-- 進捗情報 -->
    <div id="result">
        ヒットで1塁
    </div>
</body>`
	endhtml := `<body>

    <!-- イニング情報 -->
    <div class="live">
        <em>9回裏</em>
    </div>

    <!-- スコア情報 -->
    <div class="score">
        <table>
            <tr>
                <td class="nm act">オ</td>
                <td>0</td>
            </tr>
            <tr>
                <td class="nm">デ</td>
                <td>2</td>
            </tr>
        </table>
    </div>
    <!-- 進捗情報 -->
    <div id="result">
        試合終了
    </div>
</body>`

	t.Run("Success get score", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{"5回裏", "0", "2", "山田 太郎", "ヒットで1塁"},
		}
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(inninghtml))
		scoredata, err := GetMatchScore(doc)

		//エラーが発生せずに終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Equal(t, expected, scoredata)

	})
	t.Run("Success get endscore", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{"9回裏", "0", "2", "", "試合終了"},
		}
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(endhtml))
		scoredata, err := GetMatchScore(doc)

		//エラーが発生せずに終了しているか確認
		assert.NoError(t, err)

		//期待値と結果が一致しているか確認
		assert.Equal(t, expected, scoredata)

	})

}

// TestGetmatchscoreprod： 本番から試合情報取得するケースをテスト
func TestGetmatchscoreprod(t *testing.T) {
	scraper := utils.URLService{}
	t.Run("success get score", func(t *testing.T) {

		//期待値を設定
		expected := [][]string{
			{"9回表", "0", "2", "福永 奨", "二フライ\n            146km/h ストレート"},
		}

		//取得元のURLを定義
		url := "https://baseball.yahoo.co.jp/npb/game/2021030007/score?index=0910500"

		// 野球速報サイトからデータを取得
		res, err := scraper.GetURL(url)
		assert.NoError(t, err)

		doc, err := scraper.GetBody(res)
		assert.NoError(t, err)

		score, err := GetMatchScore(doc)
		assert.NoError(t, err)

		assert.Equal(t, expected, score)
	})

	t.Run("success get score to gameset", func(t *testing.T) {
		//期待値を設定
		expected := [][]string{
			{"試合終了", "0", "2", "", "試合終了\n            3回戦：DeNA 2勝0敗1分"},
		}

		//取得元のURLを定義
		url := "https://baseball.yahoo.co.jp/npb/game/2021030007/score"

		// 野球速報サイトからデータを取得
		res, err := scraper.GetURL(url)
		assert.NoError(t, err)

		doc, err := scraper.GetBody(res)
		assert.NoError(t, err)

		score, err := GetMatchScore(doc)
		assert.NoError(t, err)

		assert.Equal(t, expected, score)
	})
}
