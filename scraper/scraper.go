package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func get_url(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	return res,nil
}

func find_element_text(res *http.Response, element string) (string, error) {
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	return doc.Find(element).Text(), nil
}

func main() {
	// 野球速報サイトのURL
	url := "https://baseball.yahoo.co.jp/npb/schedule/"

	res,err := get_url(url)
	if err != nil {
		log.Fatal("Failed to get URL:",err)
	}

	element,err := find_element_text(res, ".bb-score__detail")
	if err != nil {
		log.Fatal("Failed to find element:",err)
	}

	fmt.Println(element)
}