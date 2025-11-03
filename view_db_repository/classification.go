package view_db_repository

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

func NewClassificationRepository(ctx context.Context) (*ClassificationRepository, error) {
	if classificationRepo != nil {
		return classificationRepo, nil
	}

	dsn := os.Getenv("VIEW_DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("ENV variable not set for VIEW_DB_DSN")
	}
	view_db_conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	classificationRepo = &ClassificationRepository{conn: view_db_conn}

	return classificationRepo, nil
}

func (repo *ClassificationRepository) LatestClassificationCreatedAt(ctx context.Context) (*time.Time, error) {
	var latestCreatedAt *time.Time
	query := fmt.Sprintf(`SELECT MAX(created_at) FROM %s`, CLASSIFICATION_TABLE)
	err := repo.conn.QueryRow(ctx, query).Scan(&latestCreatedAt)
	if err != nil {
		return nil, err
	}
	return latestCreatedAt, nil
}

func (repo *ClassificationRepository) BulkInsert(ctx context.Context, newClassifications []models.Classification, tx pgx.Tx) error {
	if len(newClassifications) == 0 {
		return nil
	}

	insertLiveDbQry := fmt.Sprintf(`
		INSERT INTO %s(
	id, name, arabic_name, suspended_at, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, CLASSIFICATION_TABLE)

	failed := []models.Classification{}
	for _, classification := range newClassifications {
		_, err := tx.Exec(ctx, insertLiveDbQry,
			classification.ID,
			classification.Name,
			classification.ArabicName,
			classification.SuspendedAt,
			classification.CreatedAt,
			classification.UpdatedAt,
			classification.DeletedAt,
		)
		if err != nil {
			log.Printf("Failed to insert classification ID %s: %s\n", classification.ID, err.Error())
			failed = append(failed, classification)
		}
	}

	if len(failed) > 0 {
		log.Printf("Failed to insert %d classifications.\n", len(failed))
	}

	return nil
}
