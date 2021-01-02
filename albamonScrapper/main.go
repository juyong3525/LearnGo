package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// &Page=1
var baseURL string = "http://www.albamon.com/search/Recruit?Keyword=%EB%8D%B0%EC%9D%B4%ED%84%B0&SiCode=&GuCode=&BigPartCode=&PartCode=&Gender=&Age=&Term=1&IncludeKeyword=&ExcludeKeyword=&IsExcludeDuplication=True"

type extractedJob struct {
	title    string
	location string
	salary   string
	link     string
}

func main() {
	totalPages := getPages(baseURL)
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}
func getPage(page int) {
	pageURL := baseURL + "&page=" + strconv.Itoa(page+1)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	fmt.Println("Requesting\n", pageURL)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	booth := doc.Find(".booth")
	booth.Each(func(i int, card *goquery.Selection) {
		extractJob(card)
	})
}
func extractJob(card *goquery.Selection) extractedJob {
	location := card.Find(".local").Text()
	salary := card.Find(".etc").Find("span").First().Text()
	title := card.Find("dt").First().Text()
	id, _ := card.Find("dt").Find("a").Attr("href")
	link := "http://www.albamon.com/" + id
	return extractedJob{
		location: location,
		salary:   salary,
		title:    title,
		link:     link,
	}
}
func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".listPaging").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length() + 1
	})
	return pages
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
