package main

import (
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}

	svc := service.NewRecordService()
	records, err := svc.GetNewRecords(nil)
	if err != nil {
		log.Panicln("Error fetching new records:", err)
	}

	incSvc := service.NewIncidentService()
	incidents, err := incSvc.IncidentsFromRecords(records)
	if err != nil {
		log.Panicln("Error converting records to incidents:", err)
	}

	log.Printf("Converted %d records to incidents\n", len(incidents))
	log.Printf("Fetched %d new records\n", len(records))
}
