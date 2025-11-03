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

	recordLiveDbRepo, err := live_db_repository.NewRecordRepository(ctx)
	if err != nil {
		log.Println("Failed to connect to live database: " + err.Error())
		return err
	}

	classificationLiveDbRepo, err := live_db_repository.NewClassificationRepository(ctx)
	if err != nil {
		log.Println("Failed to connect to live database: " + err.Error())
		return err
	}

	classificationViewDbRepo, err := view_db_repository.NewClassificationRepository(ctx)
	if err != nil {
		log.Println("Failed to connect to view database: " + err.Error())
		return err
	}
	recordViewDbRepo := view_db_repository.NewRecordRepository(ctx)

	// Step 1 => Getting the latest created at time among the records
	latestRecordCreatedAt, err := recordViewDbRepo.LatestRecordCreatedAt(ctx)
	if err != nil {
		log.Println("Failed to get latest record created_at: " + err.Error())
		return err
	}

	latestClassificationCreatedAt, err := classificationViewDbRepo.LatestClassificationCreatedAt(ctx)
	if err != nil {
		log.Println("Failed to get latest classification created_at: " + err.Error())
		return err
	}

	// Step 2 => Getting new data from live db
	newRecords, err := recordLiveDbRepo.GetNewRecords(ctx, latestRecordCreatedAt)
	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	newClassifications, err := classificationLiveDbRepo.GetNewClassifications(ctx, latestClassificationCreatedAt)
	if err != nil {
		log.Println("Failed to get new classifications: " + err.Error())
		return err
	}

	// Step 3 => Transaction to append new datas
	tx, err := viewDbConn.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction: " + err.Error())
		return err
	}

	if err := recordViewDbRepo.BulkInsert(ctx, newRecords, tx); err != nil {
		log.Println("Failed to sync compose records table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	if err := classificationViewDbRepo.BulkInsert(ctx, newClassifications, tx); err != nil {
		log.Println("Failed to sync compose classification table. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to sync commit transaction. Error: ", err.Error())
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
