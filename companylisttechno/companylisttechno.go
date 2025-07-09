package companylisttechno

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

type Company struct {
	ID    int
	Name  string
	Email string
	URL   string
}

func ScrapeTechnoparkCompanies() {

	base := "https://technopark.in/company-details/"
	startID := 5610
	endID := 6076

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var companies []Company

	for id := startID; id <= endID; id++ {
		url := fmt.Sprintf("%s%d", base, id)
		fmt.Printf("🔍 Scraping: %s\n", url)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("❌ Failed to parse HTML from ID %d: %v", id, err)
			continue
		}

		// Get email
		email := ""
		doc.Find("a[href^='mailto:']").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			email = strings.TrimPrefix(href, "mailto:")
		})

		// Get company name (from page title or a heading)
		companyName := strings.TrimSpace(doc.Find("title").Text())
		if companyName == "" {
			companyName = fmt.Sprintf("Company %d", id)
		}

		if email != "" {
			companies = append(companies, Company{
				ID:    id,
				Name:  companyName,
				Email: email,
				URL:   url,
			})
			fmt.Printf("✅ %d: %s | %s\n", id, companyName, email)
		}

		time.Sleep(300 * time.Millisecond)
	}

	// Write to Excel
	err := writeToExcel("technopark_company_emails.xlsx", companies)
	if err != nil {
		log.Fatalf("❌ Error writing Excel file: %v", err)
	}

	fmt.Println("✅ All done. File saved: technopark_company_emails.xlsx")
}

func writeToExcel(filename string, data []Company) error {
	f := excelize.NewFile()
	sheet := "Companies"
	index, err := f.NewSheet(sheet)
	if err != nil {
		log.Fatalf("❌ Failed to create sheet: %v", err)
	}

	headers := []string{"ID", "Company Name", "Email", "URL"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, company := range data {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), company.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), company.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), company.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), company.URL)
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filename)
}
