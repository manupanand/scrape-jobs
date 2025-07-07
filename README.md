# Job Scraping Tool in Go

This repository contains a simple job scraping tool written in Go. It is designed to scrape job listings from websites like Infopark and Technopark, extract relevant job details, and save the results in CSV format.

---

## Repository Structure

```
/infopark/infopark.go       # Scraper logic for Infopark job listings
/technopark/technopark.go   # Scraper logic for Technopark job listings
/main.go                    # Main entry point to run the scraping tool
```

---

## Getting Started

### Prerequisites

- Go 1.18+ installed on your machine
- Internet connection to access target job listing websites

### Running the Tool

1. Open your terminal or command prompt.
2. Navigate to the project root directory:

```
cd scrape-jobs
```

3. Run the scraping tool with:

```
go mod tidy
go run main.go
```

The tool will execute the scrapers and output CSV files with the scraped job data.

---

## Important Notes

- **Educational Purpose Only:** This tool is built for studying web scraping techniques and learning Go programming.  
- **Ethical Use:** Please use this tool responsibly and do not use it for phishing, spamming, or any unethical activities.  
- **Respect Website Terms:** Always respect the target websites' terms of service and robots.txt rules.

---

## Features

- Scrapes job listings from Infopark and Technopark websites
- Extracts job title, company name, email, job type, location, and job URL
- Supports keyword-based filtering for relevant job searches
- Saves results in CSV format for easy analysis

---

## License

This project is licensed under the GNU GP License



