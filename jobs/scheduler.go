package jobs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func StartScheduler() {
	// Load .env file
	_ = godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Read time configs
	hourStr := os.Getenv("SCHEDULE_HOUR")
	minuteStr := os.Getenv("SCHEDULE_MINUTE")
	timezone := os.Getenv("TIMEZONE")

	hour, _ := strconv.Atoi(hourStr)
	minute, _ := strconv.Atoi(minuteStr)

	if hourStr == "" {
		hour = 0 // default midnight
	}
	if minuteStr == "" {
		minute = 0
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil || timezone == "" {
		loc = time.Local
	}

	// Convert to CRON expression (minute hour * * *)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{Location: loc},
	)

	// Register your daily task
	_, err = scheduler.Register(cronExpr, NewDailyJobTask())
	if err != nil {
		log.Fatalf("❌ Could not register task: %v", err)
	}

	fmt.Printf("🕒 Scheduler started — will run at %02d:%02d (%s)\n", hour, minute, loc)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
