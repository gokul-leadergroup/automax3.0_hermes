package live_db_repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/models"
	"github.com/jackc/pgx/v5"
)

const CLASSIFICATION_TABLE = "classification"

var classificationRepo *ClassificationRepository

type ClassificationRepository struct {
	conn *pgx.Conn
}

func NewClassificationRepository() *ClassificationRepository {
	if classificationRepo != nil {
		return classificationRepo
	}

	live_db_dsn := os.Getenv("LIVE_DB_DSN")
	if live_db_dsn == "" {
		log.Panicln("LIVE_DB_DSN environment variable is not set")
	}

	live_db_conn, err := pgx.Connect(context.Background(), live_db_dsn)
	if err != nil {
		log.Panicln("Failed to connect to live database: " + err.Error())
	}
	return &ClassificationRepository{conn: live_db_conn}
}

func (repo *ClassificationRepository) GetNewClassifications(sinceTime *time.Time) ([]models.Classification, error) {
	var qry = ""
	if sinceTime == nil {
		qry = fmt.Sprintf(`
		SELECT id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at
		FROM %s 
		WHERE deleted_at IS NULL;
	`, CLASSIFICATION_TABLE)
	} else {
		qry = fmt.Sprintf(`
		SELECT id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at
		FROM %s 
		WHERE created_at > '%s' AND deleted_at IS NULL;
	`, CLASSIFICATION_TABLE, sinceTime.Format("2006-01-02 15:04:05-07"))
	}

	rows, err := repo.conn.Query(context.Background(), qry)
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
