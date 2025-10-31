package jobs

import (
	"context"
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/service"
	"github.com/hibiken/asynq"
)

const TaskSyncViewDbWithLiveDb = "task:sync_view_db_with_live_db"

// Task Creator
func NewDailyJobTask() *asynq.Task {
	return asynq.NewTask(TaskSyncViewDbWithLiveDb, nil)
}

// Handler
func SyncDatabases(ctx context.Context, t *asynq.Task) error {
	log.Println("✅ Running sync databases task...")

	log.Println("🔄 Syncing record table with incident table...")
	recordSvc := service.NewRecordService()
	err := recordSvc.SyncNow()
	if err != nil {
		log.Println("Failed to sync record table: " + err.Error())
		return err
		// TODO: Email notification on failure
	}

	log.Println("🔄 Syncing classification table...")
	classificationSvc := service.NewClassificationService()
	err = classificationSvc.SyncNow()
	if err != nil {
		log.Println("Failed to sync classification table: " + err.Error())
		return err
		// TODO: Email notification on failure
	}

	log.Println("✅ Sync databases task completed successfully.")

	return nil
}
