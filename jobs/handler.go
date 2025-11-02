package jobs

import (
	"context"
	"log"

	"github.com/gokul-leadergroup/automax3.0_hermes/config"
	"github.com/gokul-leadergroup/automax3.0_hermes/live_db_repository"
	"github.com/gokul-leadergroup/automax3.0_hermes/view_db_repository"
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

	recordLiveDbRepo := live_db_repository.NewRecordRepository()
	classificationLiveDbRepo := live_db_repository.NewClassificationRepository()

	classificationViewDbRepo := view_db_repository.NewClassificationRepository()
	recordViewDbRepo := view_db_repository.NewRecordRepository()

	// Step 1 => Getting the latest created at time among the records
	latestRecordCreatedAt, err := recordViewDbRepo.LatestRecordCreatedAt()
	if err != nil {
		log.Println("Failed to get latest record created_at: " + err.Error())
		return err
	}

	latestClassificationCreatedAt, err := classificationViewDbRepo.LatestClassificationCreatedAt()
	if err != nil {
		log.Println("Failed to get latest classification created_at: " + err.Error())
		return err
	}

	// Step 2 => Getting new data from live db
	newRecords, err := recordLiveDbRepo.GetNewRecords(latestRecordCreatedAt)
	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	newClassifications, err := classificationLiveDbRepo.GetNewClassifications(latestClassificationCreatedAt)
	if err != nil {
		log.Println("Failed to get new classifications: " + err.Error())
		return err
	}

	// Step 3 => Transaction to append new datas
	tx, err := viewDbConn.Begin(context.Background())
	if err != nil {
		log.Println("Failed to begin transaction: " + err.Error())
		return err
	}

	if err := recordViewDbRepo.BulkInsert(newRecords, tx); err != nil {
		log.Println("Failed to sync compose records table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	if err := classificationViewDbRepo.BulkInsert(newClassifications, tx); err != nil {
		log.Println("Failed to sync compose classification table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
