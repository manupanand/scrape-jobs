package companylist

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

type CompanyInfo struct {
	Name  string
	Email string
	Link  string
}

func ScrapeCompanyEmails() {

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var companies []CompanyInfo

	for i := 1; i <= 500; i++ {
		url := fmt.Sprintf("https://infopark.in/company-jobs/%d", i)
		fmt.Printf("🔍 Scraping: %s\n", url)

		resp, err := client.Get(url)
		if err != nil {
			log.Printf("❌ Failed to fetch URL %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("❌ Failed to parse HTML from %s: %v", url, err)
			continue
		}

		company := strings.TrimSpace(doc.Find("div.con h4").Text())
		email := ""
		doc.Find("div.con span").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if strings.Contains(text, "@") {
				email = text
			}
		})

		if company != "" && email != "" {
			companies = append(companies, CompanyInfo{
				Name:  company,
				Email: email,
				Link:  url,
			})
			fmt.Printf("✅ Found: %s | %s\n", company, email)
		} else {
			fmt.Println("⚠️ No email/company found")
		}

		time.Sleep(300 * time.Millisecond) // Be polite to the server
	}

	// Save results
	err := writeToExcel("infopark_company_emails.xlsx", companies)
	if err != nil {
		log.Fatalf("❌ Failed to save Excel: %v", err)
	}
	fmt.Println("✅ All done. Data saved to 'infopark_company_emails.xlsx'")
}

func writeToExcel(filename string, companies []CompanyInfo) error {
	f := excelize.NewFile()
	sheet := "Companies"
	
	index, err := f.NewSheet(sheet)
	if err != nil {
			log.Fatalf("❌ Failed to create sheet: %v", err)
	}

	headers := []string{"Company Name", "Email", "URL"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, c := range companies {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), c.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), c.Email)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), c.Link)
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filename)
}
