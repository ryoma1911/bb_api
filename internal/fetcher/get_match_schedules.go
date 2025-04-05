package fetcher

import (
	"baseball_report/utils"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetMatchSchedule(doc *goquery.Document) ([][]string, error) {
	todate := time.Now().Format("2006/01/02")

	// 2次元配列で試合情報を格納
	var matchData [][]string

	// 各リーグのスコア要素を取得
	utils.GetElement(doc, ".bb-score").Each(func(index int, param *goquery.Selection) {
		//試合がない場合はno card todayを返す
		if utils.GetText(param, ".bb-noData") == "" {

			// リーグ名を取得
			header := utils.GetText(param, ".bb-score__title")

			// リーグ内の各試合情報を取得
			utils.GetElement(param, ".bb-score__item").Each(func(count int, card *goquery.Selection) {
				home := utils.GetText(card, "[class*='bb-score__homeLogo']")
				away := utils.GetText(card, "[class*='bb-score__awayLogo']")
				stadium := utils.GetText(card, ".bb-score__venue")
				status := utils.GetText(card, ".bb-score__link")
				starttime := utils.GetText(card, ".bb-score__status")
				link, err := utils.GetElement(card, ".bb-score__content").Attr("href")
				if !err {
					log.Println("Link not found for the match.")
					return
				}
				link = strings.Replace(link, "index", "score", 1)

				matchData = append(matchData, []string{todate, home, away, stadium, status, starttime, link, header})
			})

		} else {
			log.Println("No card today.")
		}
	})
	return matchData, nil
}
