package jobs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
)

func StartScheduler() {
	redisAddr := os.Getenv("REDIS_ADDR")

	hourStr := os.Getenv("SCHEDULE_HOUR")
	minuteStr := os.Getenv("SCHEDULE_MINUTE")
	timezone := os.Getenv("TIMEZONE")

	hour, _ := strconv.Atoi(hourStr)
	minute, _ := strconv.Atoi(minuteStr)

	if hourStr == "" {
		hour = 0
	}
	if minuteStr == "" {
		minute = 0
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil || timezone == "" {
		loc = time.Local
	}

	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{Location: loc},
	)

	_, err = scheduler.Register(cronExpr, NewDailyJobTask())
	if err != nil {
		log.Fatalf("❌ Could not register task: %v", err)
	}

	fmt.Printf("🕒 Hermes started — will sync databases at %02d:%02d (%s)\n", hour, minute, loc)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
