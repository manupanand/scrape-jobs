package main

import (
	"github.com/manupanand/scrape-jobs/companylist"
	"github.com/manupanand/scrape-jobs/infopark"
	"github.com/manupanand/scrape-jobs/technopark"
)
func main (){
	
	technopark.TechnoParkJobs()
	infopark.InfoParkJobs()
	companylist.ScrapeCompanyEmails()
}