package jobs

import (
	"context"
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/config"
	"github.com/gokul-leadergroup/automax3.0_hermes/live_db_repository"
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

	viewDbConn, err := config.ViewDbConn()
	if err != nil {
		log.Println("Failed to connect to view database: " + err.Error())
		return err
	}

	recordRepo := live_db_repository.NewRecordRepository()
	classificationRepo := live_db_repository.NewClassificationRepository()

	latestRecordCreatedAt, err := recordRepo.LatestRecordCreatedAt()
	if err != nil {
		log.Println("Failed to get latest record created_at: " + err.Error())
		return err
	}

	latestClassificationCreatedAt, err := classificationRepo.LatestClassificationCreatedAt()
	if err != nil {
		log.Println("Failed to get latest classification created_at: " + err.Error())
		return err
	}

	newRecords, err := recordRepo.GetNewRecords(latestRecordCreatedAt)
	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	newClassifications, err := classificationRepo.GetNewClassifications(latestClassificationCreatedAt)
	if err != nil {
		log.Println("Failed to get new classifications: " + err.Error())
		return err
	}

	// Transaction to insert records and classifications
	tx, err := viewDbConn.Begin(context.Background())
	if err != nil {
		log.Println("Failed to begin transaction: " + err.Error())
		return err
	}

	if err := recordRepo.SyncViewDbWithLiveDB(newRecords, tx); err != nil {
		log.Println("Failed to sync compose records table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	if err := classificationRepo.SyncViewDbWithLiveDB(newClassifications, tx); err != nil {
		log.Println("Failed to sync compose classification table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
