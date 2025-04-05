package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// URLHandler インターフェースで関数を抽象化
type URLHandler interface {
	GetURL(url string) (*http.Response, error)
	GetBody(res *http.Response) (*goquery.Document, error)
}

// URLService デフォルトの実装
type URLService struct{}

// サイトのURLからHTTPレスポンスを取得
func (u *URLService) GetURL(url string) (*http.Response, error) {
	//URLが空の場合はエラーを返す
	if url == "" {
		return nil, fmt.Errorf("response body is empty")
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return res, nil
}

func (u *URLService) GetBody(res *http.Response) (*goquery.Document, error) {
	defer res.Body.Close() // リソースリーク防止

	// レスポンスボディの内容をチェック
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("response body is empty")
	}

	// goquery ドキュメントを作成
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}
	return doc, nil
}

// GetElement は *goquery.Document または *goquery.Selection を受け取り、指定された要素を取得する
func GetElement(body interface{}, element string) *goquery.Selection {
	var selection *goquery.Selection

	// 型アサーションで *goquery.Document か *goquery.Selection を判定
	switch v := body.(type) {
	case *goquery.Document:
		selection = v.Find(element)
	case *goquery.Selection:
		selection = v.Find(element)
	default:
		return nil // 無効な型の場合は nil を返す
	}

	return selection
}

// GetText は *goquery.Document または *goquery.Selection を受け取り、指定された要素のテキストを取得する
func GetText(body interface{}, element string) string {
	var selection *goquery.Selection

	// 型アサーションで *goquery.Document か *goquery.Selection を判定
	switch v := body.(type) {
	case *goquery.Document:
		selection = v.Find(element)
	case *goquery.Selection:
		selection = v.Find(element)
	default:
		return "" // 無効な型の場合は空文字を返す
	}

	return strings.TrimSpace(selection.Text())
}
