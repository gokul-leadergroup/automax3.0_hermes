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
	records, err := svc.GetNewIncidents(nil)
	if err != nil {
		log.Panicln("Error fetching new incidents:", err)
	}

	log.Printf("Fetched %d new records\n", len(records))
}
