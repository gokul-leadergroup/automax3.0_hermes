package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/service"
	"github.com/hibiken/asynq"
)

const TaskSyncViewDbWithLiveDb = "task:sync_view_db_with_live_db"

// No payload → just create an empty task
func NewDailyJobTask() *asynq.Task {
	return asynq.NewTask(TaskSyncViewDbWithLiveDb, nil)
}

// Handler: your business logic here
func SyncDatabases(ctx context.Context, t *asynq.Task) error {
	log.Println("✅ Running daily scheduled task...")
	recordSvc := service.NewRecordService()
	err := recordSvc.SyncNow()
	if err != nil {
		return err
	}

	fmt.Println("✅ Daily scheduled task completed successfully.")

	return nil
}
