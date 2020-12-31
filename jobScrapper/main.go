package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

// 기본 URL에서 각각의 페이지 URL을 만드는 함수
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL) // 각각의 page URL을 get한다
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // defer: 함수가 끝난 시점(마지막)에 실행

	doc, err := goquery.NewDocumentFromReader(res.Body) // res의 html을 가져온다
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")                      // 클래스의 값을 가져오는 메소드
	title := cleanString(card.Find(".title>a").Text()) // Text(): 클래스에 있는 string을 가져오는 메소드
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	return extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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
