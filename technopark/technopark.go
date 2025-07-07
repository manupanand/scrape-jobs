// package technopark

// import (
//     "crypto/tls"
//     "encoding/csv"
//     "fmt"
//     "net/http"
//     "os"
//     "strings"

//     "github.com/PuerkitoBio/goquery"
// )

// type Job struct {
//     ID      int
//     Title   string
//     Company string
//     Email   string
//     URL     string
// }

// func TechnoParkJobs() {
//     base := "https://technopark.in/job-details/"
//     startID := 20716
//     endID := 21599 // increase as needed

//     keywords := []string{"python", "golang", "mern", "devops","full stack"}

//     client := &http.Client{
//         Transport: &http.Transport{
//             TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//         },
//     }

//     var jobs []Job

//     for id := startID; id <= endID; id++ {
//         url := fmt.Sprintf("%s%d", base, id)
//         fmt.Printf("Checking ID: %d\n", id)

//         req, _ := http.NewRequest("GET", url, nil)
//         req.Header.Set("User-Agent", "Mozilla/5.0")

//         resp, err := client.Do(req)
//         if err != nil || resp.StatusCode != 200 {
//             continue
//         }
//         defer resp.Body.Close()

//         doc, err := goquery.NewDocumentFromReader(resp.Body)
//         if err != nil {
//             continue
//         }

//         jobTitle := strings.TrimSpace(doc.Find("div.mx-4.mt-5 h1").First().Text())
//         lowerTitle := strings.ToLower(jobTitle)

//         match := false
//         for _, keyword := range keywords {
//             if strings.Contains(lowerTitle, keyword) {
//                 match = true
//                 break
//             }
//         }
//         if !match {
//             continue // skip non-matching jobs
//         }

//         company := strings.TrimSpace(doc.Find("a[href^='/company-details/']").First().Text())

//         email := ""
//         doc.Find("a[href^='mailto:']").Each(func(i int, s *goquery.Selection) {
//             if href, exists := s.Attr("href"); exists && strings.HasPrefix(href, "mailto:") {
//                 email = strings.TrimSpace(strings.TrimPrefix(href, "mailto:"))
//             }
//         })

//         fmt.Printf("✅ ID: %d\nTitle: %s\nCompany: %s\nEmail: %s\nURL: %s\n\n",
//             id, jobTitle, company, email, url)

//         jobs = append(jobs, Job{
//             ID:      id,
//             Title:   jobTitle,
//             Company: company,
//             Email:   email,
//             URL:     url,
//         })
//     }

//     // Create or truncate file.csv for writing
//     file, err := os.Create("techno_park_jobs.csv")
//     if err != nil {
//         fmt.Println("Error creating CSV file:", err)
//         return
//     }
//     defer file.Close()

//     w := csv.NewWriter(file)
//     defer w.Flush()

//     // Write CSV header
//     if err := w.Write([]string{"ID", "Title", "Company", "Email", "URL"}); err != nil {
//         fmt.Println("Error writing CSV header:", err)
//         return
//     }

//     // Write job records
//     for _, job := range jobs {
//         record := []string{
//             fmt.Sprintf("%d", job.ID),
//             job.Title,
//             job.Company,
//             job.Email,
//             job.URL,
//         }
//         if err := w.Write(record); err != nil {
//             fmt.Println("Error writing CSV record:", err)
//             return
//         }
//     }

//	    fmt.Println("✅ CSV file 'file.csv' written successfully.")
//	}
package technopark

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
)

type Job struct {
	ID      int
	Title   string
	Company string
	Email   string
	URL     string
}

func TechnoParkJobs() {
	base := "https://technopark.in/job-details/"
	startID := 20716
	endID := 21599

	keywords := []string{"python", "golang", "mern", "devops", "full stack"}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var jobs []Job

	for id := startID; id <= endID; id++ {
		url := fmt.Sprintf("%s%d", base, id)
		fmt.Printf("Checking ID: %d\n", id)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}

		jobTitle := strings.TrimSpace(doc.Find("div.mx-4.mt-5 h1").First().Text())
		lowerTitle := strings.ToLower(jobTitle)

		match := false
		for _, keyword := range keywords {
			if strings.Contains(lowerTitle, keyword) {
				match = true
				break
			}
		}
		if !match {
			continue
		}

		company := strings.TrimSpace(doc.Find("a[href^='/company-details/']").First().Text())

		email := ""
		doc.Find("a[href^='mailto:']").Each(func(i int, s *goquery.Selection) {
			if href, exists := s.Attr("href"); exists && strings.HasPrefix(href, "mailto:") {
				email = strings.TrimSpace(strings.TrimPrefix(href, "mailto:"))
			}
		})

		fmt.Printf("✅ ID: %d\nTitle: %s\nCompany: %s\nEmail: %s\nURL: %s\n\n",
			id, jobTitle, company, email, url)

		jobs = append(jobs, Job{
			ID:      id,
			Title:   jobTitle,
			Company: company,
			Email:   email,
			URL:     url,
		})
	}

	// Write to Excel (.xlsx)
	f := excelize.NewFile()
	sheet := "Jobs"
	index, err := f.NewSheet(sheet)
	if err != nil {
		log.Fatalf("❌ Failed to create sheet: %v", err)
	}

	// Write header
	headers := []string{"ID", "Title", "Company", "Email", "URL"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Write rows
	for rowIndex, job := range jobs {
		values := []interface{}{job.ID, job.Title, job.Company, job.Email, job.URL}
		for colIndex, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	f.SetActiveSheet(index)
	err = f.SaveAs("techno_park_jobs.xlsx")
	if err != nil {
		fmt.Println("❌ Error saving Excel file:", err)
		return
	}

	fmt.Println("✅ Excel file 'techno_park_jobs.xlsx' written successfully.")
}
