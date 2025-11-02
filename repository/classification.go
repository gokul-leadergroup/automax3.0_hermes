package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/models"
	"github.com/jackc/pgx/v5"
)

const LIVE_DB_CLASSIFICATION_TABLE = "classification"
const VIEW_DB_CLASSIFICATION_TABLE = "classification"

var classificationRepo *ClassificationRepository

type ClassificationRepository struct {
	liveDb *pgx.Conn
	viewDb *pgx.Conn
}

func NewClassificationRepository() *ClassificationRepository {
	if classificationRepo != nil {
		return classificationRepo
	}

	live_db_dsn := os.Getenv("LIVE_DB_DSN")
	if live_db_dsn == "" {
		log.Panicln("LIVE_DB_DSN environment variable is not set")
	}

	live_db_conn, err := PgConx(live_db_dsn)
	if err != nil {
		log.Panicln("Failed to connect to live database: " + err.Error())
	}

	view_db_dsn := os.Getenv("VIEW_DB_DSN")
	var view_db_conn *pgx.Conn = nil
	if view_db_dsn != "" {
		view_db_conn, err = PgConx(view_db_dsn)
		if err != nil {
			log.Panicln("Failed to connect to view database: " + err.Error())
		}
	}
	return &ClassificationRepository{liveDb: live_db_conn, viewDb: view_db_conn}
}

func (repo *ClassificationRepository) LatestClassificationCreatedAt() (*time.Time, error) {
	var latestCreatedAt *time.Time
	query := fmt.Sprintf(`SELECT MAX(created_at) FROM %s`, VIEW_DB_CLASSIFICATION_TABLE)
	err := repo.viewDb.QueryRow(context.Background(), query).Scan(&latestCreatedAt)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	return latestCreatedAt, nil
}

func (repo *ClassificationRepository) GetNewClassifications(sinceTime *time.Time) ([]models.Classification, error) {
	var qry = ""
	if sinceTime == nil {
		qry = fmt.Sprintf(`
		SELECT id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at
		FROM %s 
		WHERE deleted_at IS NULL;
	`, LIVE_DB_CLASSIFICATION_TABLE)
	} else {
		qry = fmt.Sprintf(`
		SELECT id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at
		FROM %s 
		WHERE created_at > '%s' AND deleted_at IS NULL;
	`, LIVE_DB_CLASSIFICATION_TABLE, sinceTime.Format("2006-01-02 15:04:05-07"))
	}

	rows, err := repo.liveDb.Query(context.Background(), qry)
	if err != nil {
		log.Println("Failed to execute query: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var classifications []models.Classification
	for rows.Next() {
		var classification models.Classification
		if err := rows.Scan(&classification.ID, &classification.Name, &classification.ArabicName, &classification.SuspendedAt, &classification.CreatedAt, &classification.UpdatedAt, &classification.DeletedAt); err != nil {
			log.Println("Failed to scan row: " + err.Error())
			return nil, err
		}
		classifications = append(classifications, classification)
	}
	return classifications, nil
}

func (repo *ClassificationRepository) SyncViewDbWithLiveDB(classifications []models.Classification) error {
	if len(classifications) == 0 {
		log.Println("No new records to sync.")
		return nil
	}

	log.Printf("Found %d new records to sync.\n", len(classifications))

	insertLiveDbQry := fmt.Sprintf(`
		INSERT INTO %s(
	id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, VIEW_DB_CLASSIFICATION_TABLE)

	failed := []models.Classification{}
	for _, classification := range classifications {
		_, err := repo.viewDb.Exec(context.Background(), insertLiveDbQry,
			classification.ID,
			classification.Name,
			classification.ArabicName,
			classification.SuspendedAt,
			classification.CreatedAt,
			classification.UpdatedAt,
			classification.DeletedAt,
		)
		if err != nil {
			log.Println("Failed to insert classification: " + err.Error())
			failed = append(failed, classification)
		}
	}

	if len(failed) > 0 {
		log.Printf("Failed to sync %d classifications.\n", len(failed))
	} else {
		log.Println("All classifications synced successfully.")
	}

	return nil
}

func (repo *ClassificationRepository) SyncNow() error {
	latestCreatedAt, err := repo.LatestClassificationCreatedAt()
	if err != nil {
		log.Println("Failed to get latest created_at: " + err.Error())
		return err
	}

	if latestCreatedAt == nil {
		log.Println("Table is empty — no records found.")
	} else {
		log.Println("Latest created_at:", latestCreatedAt)
	}

	var classifications []models.Classification
	classifications, err = repo.GetNewClassifications(latestCreatedAt)

	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	return repo.SyncViewDbWithLiveDB(classifications)
}
