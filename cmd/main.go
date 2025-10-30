package main

import (
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/jobs"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TaskSyncViewDbWithLiveDb, jobs.SyncDatabases)
}
