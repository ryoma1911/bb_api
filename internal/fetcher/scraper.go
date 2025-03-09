package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// サイトのURLからHTTPレスポンスを取得
func getURL(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return res, nil
}

// goqueryドキュメントを作成
func getBody(res *http.Response) (*goquery.Document, error) {
	defer res.Body.Close() // リソースリーク防止
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}
	return doc, nil
}

func main() {
	// 野球速報サイトのURL
	url := "https://baseball.yahoo.co.jp/npb/schedule/"

	// URLからレスポンスを取得
	res, err := getURL(url)
	if err != nil {
		log.Fatal("Failed to get URL:", err)
	}

	// レスポンスからドキュメントを取得
	doc, err := getBody(res)
	if err != nil {
		log.Fatal("Failed to get body:", err)
	}

	// 各リーグのスコア要素を取得
	doc.Find(".bb-score").Each(func(index int, param *goquery.Selection) {
		// リーグ名を取得
		header := strings.TrimSpace(param.Find(".bb-score__title").Text())
		fmt.Println("league:", header)

		// リーグ内の各試合情報を取得
		param.Find(".bb-score__item").Each(func(count int, card *goquery.Selection) {
			home := strings.TrimSpace(card.Find("[class*='bb-score__homeLogo']").Text())
			away := strings.TrimSpace(card.Find("[class*='bb-score__awayLogo']").Text())
			stadium := strings.TrimSpace(card.Find(".bb-score__venue").Text())
			status := strings.TrimSpace(card.Find(".bb-score__status").Text())
			starttime := strings.TrimSpace(card.Find(".bb-score__link").Text())
			fmt.Println("home:", home)
			fmt.Println("away:", away)
			fmt.Println("stadium:", stadium)
			fmt.Println("status:", status)
			fmt.Println("start:", starttime)
		})
	})
}
