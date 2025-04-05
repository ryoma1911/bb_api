package utils

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// HTMLのモックデータ
const mockHTML = `
<html>
<head><title>Test Page</title></head>
<body>
    <h1 id="title">Hello World</h1>
    <p class="message">This is a test.</p>
    <div class="container">
        <span class="info">Some info</span>
    </div>
</body>
</html>
`

// GetURL のテスト
func TestGetURL(t *testing.T) {
	scraper := URLService{}

	// モックサーバーを作成して、HTTPレスポンスをシミュレート
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // ステータスコード200を返す
		w.Write([]byte("<html><body><h1>Hello, World!</h1></body></html>"))
	}))
	defer mockServer.Close()

	t.Run("Valid URL", func(t *testing.T) {
		// モックサーバーが正しく動作するか
		res, err := scraper.GetURL(mockServer.URL)
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		// 無効なURLでエラーが発生すること
		_, err := scraper.GetURL("http://invalid-url")
		assert.Error(t, err)
	})

	t.Run("Empty URL", func(t *testing.T) {
		// 空のURLでエラーが発生すること
		_, err := scraper.GetURL("")
		assert.Error(t, err)
	})
}

type badReader struct{}

func (b *badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

// GetBody のテスト
func TestGetBody(t *testing.T) {
	scraper := URLService{}
	// モックレスポンスを作成
	mockRes := httptest.NewRecorder()
	mockRes.WriteString("<html><body><h1>Hello, World!</h1></body></html>")
	res := mockRes.Result()

	t.Run("Valid HTML", func(t *testing.T) {
		// 有効なHTMLレスポンスを渡して、goquery.Document を取得する
		doc, err := scraper.GetBody(res)
		assert.NoError(t, err)

		// 正しい内容が返っているか確認
		assert.Equal(t, doc.Find("h1").Text(), "Hello, World!")
	})

	t.Run("Empty HTTP response", func(t *testing.T) {
		// 空のレスポンスを作成
		emptyRes := httptest.NewRecorder()
		emptyResponse := emptyRes.Result()

		// GetBody を実行
		_, err := scraper.GetBody(emptyResponse)

		// エラーが発生することを確認
		assert.Error(t, err)
	})

	t.Run("Invalid HTML response", func(t *testing.T) {
		// エラーを返す `io.Reader`
		invalidResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(&badReader{}),
		}

		// GetBody を実行
		_, err := scraper.GetBody(invalidResponse)

		// goquery は軽微なHTMLエラーを許容するので、明確なエラーが出るか確認
		assert.Error(t, err)
	})
}

// goquery.Document を取得するヘルパー関数
func getTestDocument() *goquery.Document {
	res := httptest.NewRecorder()
	res.WriteString(mockHTML)
	response := res.Result()
	defer response.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(response.Body)
	return doc
}

// TestGetElement: *goquery.Document または *goquery.Selection を受け取るケースをテスト
func TestGetElement(t *testing.T) {
	doc := getTestDocument()

	t.Run("Find h1 element from *goquery.Document", func(t *testing.T) {
		// *goquery.Document から <h1> 要素を取得する
		selection := GetElement(doc, "h1")
		assert.NotNil(t, selection)
		assert.Greater(t, selection.Length(), 0)
	})

	t.Run("Find non-existent element", func(t *testing.T) {
		// 存在しない要素 (footer) を検索し、空の selection を返すことを確認
		selection := GetElement(doc, "footer")
		assert.NotNil(t, selection)
		assert.Equal(t, selection.Length(), 0)
	})

	t.Run("Find element from *goquery.Selection", func(t *testing.T) {
		// *goquery.Selection (container) から .info 要素を取得する
		container := GetElement(doc, ".container")
		selection := GetElement(container, ".info")
		assert.NotEqual(t, selection, nil)
	})

	t.Run("Invalid input type", func(t *testing.T) {
		// 無効な型 (string) を渡した場合、nil を返すことを確認
		selection := GetElement("invalid", "h1")
		assert.Nil(t, selection)
	})
}

// TestGetText: *goquery.Document または *goquery.Selection からテキストを取得するケースをテスト
func TestGetText(t *testing.T) {
	doc := getTestDocument()

	t.Run("Get text from h1", func(t *testing.T) {
		// *goquery.Document から <h1> のテキストを取得する
		text := GetText(doc, "h1")
		expected := "Hello World"
		assert.Equal(t, text, expected)

	})

	t.Run("Get text from p", func(t *testing.T) {
		// *goquery.Document から <p class="message"> のテキストを取得する
		text := GetText(doc, ".message")
		expected := "This is a test."
		assert.Equal(t, text, expected)
	})

	t.Run("Get text from *goquery.Selection", func(t *testing.T) {
		// *goquery.Selection (container) を使って .info のテキストを取得する
		selection := GetElement(doc, ".container")
		text := GetText(selection, ".info")
		expected := "Some info"
		assert.Equal(t, text, expected)
	})

	t.Run("Get text from non-existent element", func(t *testing.T) {
		// 存在しない要素を指定した場合、空文字が返ることを確認
		text := GetText(doc, ".does-not-exist")
		assert.Equal(t, text, "")
	})

	t.Run("Invalid input type", func(t *testing.T) {
		// 無効な型 (string) を渡した場合、空文字が返ることを確認
		text := GetText("invalid", "h1")
		assert.Equal(t, text, "")
	})
}
