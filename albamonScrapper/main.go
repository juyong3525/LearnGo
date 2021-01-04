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

var baseURL string = "http://www.albamon.com/search/Recruit?Keyword=%EB%8D%B0%EC%9D%B4%ED%84%B0&SiCode=&GuCode=&BigPartCode=&PartCode=&Gender=&Age=&Term=1&IncludeKeyword=&ExcludeKeyword=&IsExcludeDuplication=True"

type extractedJob struct {
	title    string
	location string
	salary   string
	link     string
}

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
}

func getPage(page int, mainC chan []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := baseURL + "&page=" + strconv.Itoa(page+1)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	fmt.Println("Requesting\n", pageURL)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	booths := doc.Find(".booth")
	booths.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < booths.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan extractedJob) {
	location := CleanString(card.Find(".local").Text())
	salary := CleanString(card.Find(".etc").Find("span").First().Text())
	title := CleanString(card.Find("dt").First().Text())
	id, _ := card.Find("dt").Find("a").Attr("href")
	link := "http://www.albamon.com/" + id

	c <- extractedJob{
		location: location,
		salary:   salary,
		title:    title,
		link:     link,
	}
}

// CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"title", "location", "salary", "link"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.title, job.location, job.salary, job.link}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
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
