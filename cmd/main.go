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
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 5},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TaskSyncViewDbWithLiveDb, jobs.SyncDatabases)

	go func() {
		jobs.StartScheduler()
	}()

	fmt.Println("🚀 Asynq worker and scheduler started...")
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
