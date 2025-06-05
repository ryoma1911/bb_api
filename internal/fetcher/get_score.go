package fetcher

import (
	"baseball_report/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetMatchScore(doc *goquery.Document) ([][]string, error) {

	// 2次元配列で試合進捗を格納
	var scoreData [][]string
	var teamscore []string

	//イニングを取得
	inning := utils.GetText(doc, ".live em")
	//両チームのスコアを取得
	utils.GetElement(doc, "tr").Each(func(i int, s *goquery.Selection) {
		score := s.Find("td").Eq(1).Text()
		if score != "" {
			teamscore = append(teamscore, score)
		}
	})
	//進捗を取得
	result := utils.GetText(doc, "div#result")

	// 打者情報を取得
	batter := utils.GetText(doc, "table#batt a")

	scoreData = append(scoreData, []string{inning, teamscore[1], teamscore[0], batter, result})

	return scoreData, nil
}
