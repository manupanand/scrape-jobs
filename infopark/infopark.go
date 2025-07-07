// package infopark

// import (
// 	"crypto/tls"
// 	"encoding/csv"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/PuerkitoBio/goquery"
// )

// type Job struct {
// 	Title   string
// 	Company string
// 	Email   string
// 	Link    string
// }

// func InfoParkJobs() {
// 	client := &http.Client{
// 		Timeout: 15 * time.Second,
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 		},
// 	}

// 	keywords := []string{"devops", "python", "go", "full stack"}

// 	var jobs []Job

// 	for _, keyword := range keywords {
// 		fmt.Printf("Searching jobs for keyword: %s\n", keyword)
// 		page := 1

// 		for {
// 			url := fmt.Sprintf("https://infopark.in/companies/job-search?search=%s&page=%d", strings.ReplaceAll(keyword, " ", "+"), page)
// 			resp, err := client.Get(url)
// 			if err != nil {
// 				log.Printf("Failed to fetch page %d for keyword %s: %v", page, keyword, err)
// 				break
// 			}
// 			defer resp.Body.Close()

// 			doc, err := goquery.NewDocumentFromReader(resp.Body)
// 			if err != nil {
// 				log.Printf("Failed to parse page %d for keyword %s: %v", page, keyword, err)
// 				break
// 			}

// 			jobsOnPage := 0

// 			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
// 				title := strings.TrimSpace(s.Find("td.head").Text())
// 				company := strings.TrimSpace(s.Find("td.date").Text())
// 				link, exists := s.Find("a").Attr("href")

// 				// Only add if title contains the keyword (case-insensitive)
// 				if exists && strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
// 					if !strings.HasPrefix(link, "http") {
// 						link = "https://infopark.in" + link
// 					}
// 					email := fetchEmail(client, link)
// 					if email == "" {
// 						email = "Not found"
// 					}
// 					jobs = append(jobs, Job{
// 						Title:   title,
// 						Company: company,
// 						Email:   email,
// 						Link:    link,
// 					})
// 					jobsOnPage++
// 				}
// 			})

// 			if jobsOnPage == 0 {
// 				break // no more jobs on this page for this keyword
// 			}
// 			fmt.Printf("Scraped Page %d for keyword '%s': %d job(s)\n", page, keyword, jobsOnPage)
// 			page++
// 			// To avoid hammering the server, you might add a short sleep here
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}

// 	fmt.Printf("\nTotal Jobs Found: %d\n", len(jobs))
// 	for _, job := range jobs {
// 		fmt.Printf("%s | %s | %s | %s\n", job.Title, job.Company, job.Email, job.Link)
// 	}

// 	err := writeJobsToCSV("info_park_jobs.csv", jobs)
// 	if err != nil {
// 		log.Fatalf("Failed to write CSV file: %v", err)
// 	}
// 	fmt.Println("✅ CSV file 'jobs.csv' written successfully.")
// }

// func fetchEmail(client *http.Client, url string) string {
// 	resp, err := client.Get(url)
// 	if err != nil {
// 		log.Printf("Failed to fetch job detail page: %v", err)
// 		return ""
// 	}
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		log.Printf("Failed to parse job detail page: %v", err)
// 		return ""
// 	}

// 	var email string
// 	doc.Find("div.con span").Each(func(i int, s *goquery.Selection) {
// 		text := strings.TrimSpace(s.Text())
// 		if strings.Contains(text, "@") {
// 			email = text
// 		}
// 	})

// 	return email
// }

// func writeJobsToCSV(filename string, jobs []Job) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	if err := writer.Write([]string{"Title", "Company", "Email", "Link"}); err != nil {
// 		return err
// 	}

// 	for _, job := range jobs {
// 		record := []string{job.Title, job.Company, job.Email, job.Link}
// 		if err := writer.Write(record); err != nil {
// 			return err
// 		}
// 	}

//		return nil
//	}
package infopark

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
)

type Job struct {
	Title   string
	Company string
	Email   string
	Link    string
}

func InfoParkJobs() {
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	keywords := []string{"devops", "python", "go", "full stack"}

	var jobs []Job

	for _, keyword := range keywords {
		fmt.Printf("🔍 Searching jobs for keyword: %s\n", keyword)
		page := 1

		for {
			url := fmt.Sprintf("https://infopark.in/companies/job-search?search=%s&page=%d", strings.ReplaceAll(keyword, " ", "+"), page)
			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Failed to fetch page %d for keyword %s: %v", page, keyword, err)
				break
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("Failed to parse page %d for keyword %s: %v", page, keyword, err)
				break
			}

			jobsOnPage := 0

			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				title := strings.TrimSpace(s.Find("td.head").Text())
				company := strings.TrimSpace(s.Find("td.date").Text())
				link, exists := s.Find("a").Attr("href")

				if exists && strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
					if !strings.HasPrefix(link, "http") {
						link = "https://infopark.in" + link
					}
					email := fetchEmail(client, link)
					if email == "" {
						email = "Not found"
					}
					jobs = append(jobs, Job{
						Title:   title,
						Company: company,
						Email:   email,
						Link:    link,
					})
					jobsOnPage++
				}
			})

			if jobsOnPage == 0 {
				break
			}
			fmt.Printf("✅ Scraped Page %d for keyword '%s': %d job(s)\n", page, keyword, jobsOnPage)
			page++
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("\n📋 Total Jobs Found: %d\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("%s | %s | %s | %s\n", job.Title, job.Company, job.Email, job.Link)
	}

	err := writeJobsToExcel("info_park_jobs.xlsx", jobs)
	if err != nil {
		log.Fatalf("❌ Failed to write Excel file: %v", err)
	}
	fmt.Println("✅ Excel file 'info_park_jobs.xlsx' written successfully.")
}

func fetchEmail(client *http.Client, url string) string {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to fetch job detail page: %v", err)
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse job detail page: %v", err)
		return ""
	}

	var email string
	doc.Find("div.con span").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "@") {
			email = text
		}
	})

	return email
}

func writeJobsToExcel(filename string, jobs []Job) error {
	f := excelize.NewFile()
	sheet := "Jobs"
	index, err := f.NewSheet(sheet)
	if err != nil {
			log.Fatalf("❌ Failed to create sheet: %v", err)
	}

	// Write header row
	headers := []string{"Title", "Company", "Email", "Link"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Write job data
	for rowIdx, job := range jobs {
		values := []interface{}{job.Title, job.Company, job.Email, job.Link}
		for colIdx, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filename)
}
