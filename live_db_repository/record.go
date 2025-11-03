package live_db_repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/models"
	"github.com/jackc/pgx/v5"
)

const RECORD_TABLE = "compose_record_hab_ims"

var recordRepo *RecordRepository

type RecordRepository struct {
	conn *pgx.Conn
}

func NewRecordRepository(ctx context.Context) (*RecordRepository, error) {
	if recordRepo != nil {
		return recordRepo, nil
	}

	live_db_dsn := os.Getenv("LIVE_DB_DSN")
	if live_db_dsn == "" {
		return nil, fmt.Errorf("LIVE_DB_DSN environment variable is not set")
	}

	live_db_conn, err := pgx.Connect(ctx, live_db_dsn)
	if err != nil {
		return nil, err
	}

	return &RecordRepository{conn: live_db_conn}, nil
}

func (repo *RecordRepository) GetNewRecords(ctx context.Context, sinceTime *time.Time) ([]models.Record, error) {
	var qry = ""
	if sinceTime == nil {
		qry = fmt.Sprintf(`
		SELECT id, revision, rel_module, "values", meta, rel_namespace, 
		       created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by 
		FROM %s 
		WHERE deleted_at IS NULL;
	`, RECORD_TABLE)
	} else {
		qry = fmt.Sprintf(`
		SELECT id, revision, rel_module, "values", meta, rel_namespace, 
		       created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by 
		FROM %s 
		WHERE created_at > '%s' AND deleted_at IS NULL;
	`, RECORD_TABLE, sinceTime.Format("2006-01-02 15:04:05-07"))
	}

	rows, err := repo.conn.Query(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var record models.Record
		if err := rows.Scan(&record.ID, &record.Revision, &record.ModuleID, &record.Values, &record.Meta, &record.NamespaceID, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt, &record.OwnedBy, &record.CreatedBy, &record.UpdatedBy, &record.DeletedBy); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
