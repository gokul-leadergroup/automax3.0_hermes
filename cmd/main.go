package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gokul-leadergroup/automax3.0_hermes/jobs"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panicln("Error loading .env file:", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 5},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TaskSyncViewDbWithLiveDb, jobs.SyncDatabases)

	// Start scheduler in a separate goroutine
	go jobs.StartScheduler()

	fmt.Println("🚀 Asynq worker and scheduler started...")
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
