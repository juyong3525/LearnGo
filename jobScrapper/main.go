package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

// 기본 URL
var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

// 기본 URL에서 각각의 페이지 URL을 만드는 함수
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)

	res, err := http.Get(pageURL) // 각각의 page URL을 get한다
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // defer: 함수가 끝난 시점(마지막)에 실행

	doc, err := goquery.NewDocumentFromReader(res.Body) // res의 html을 가져온다
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")         // 클래스의 값을 가져오는 메소드
		title := card.Find(".title>a").Text() // Text(): 클래스에 있는 string을 가져오는 메소드
		location := card.Find(".sjcl").Find(".location").Text()

		fmt.Println(id, title, location)
	})
}

// 전체 페이지 개수를 찾는 함수
func getPages() int {
	pages := 0 // 전체 페이지 개수를 0으로 초기화
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // defer: 함수가 끝난 시점(마지막)에 실행

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// https://github.com/PuerkitoBio/goquery 참고
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length() // a href의 총 개수를 확인(현재 페이지를 제외한 다음 페이지 전체 링크 개수를 알아냄)
	})

	return pages
}

// 에러를 체크하는 함수
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// status 코드를 체크하는 함수
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
