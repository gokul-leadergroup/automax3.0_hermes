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

const LIVE_DB_RECORD_TABLE = "compose_record_hab_ims"
const VIEW_DB_INCIDENT_TABLE = "compose_record_hab_ims"

var recordRepo *RecordRepository

type RecordRepository struct {
	liveDb *pgx.Conn
	viewDb *pgx.Conn
}

func NewRecordRepository() *RecordRepository {
	if recordRepo != nil {
		return recordRepo
	}

	live_db_dsn := os.Getenv("LIVE_DB_DSN")
	if live_db_dsn == "" {
		log.Panicln("LIVE_DB_DSN environment variable is not set")
	}

	live_db_conn, err := connectToDB(live_db_dsn)
	if err != nil {
		log.Panicln("Failed to connect to live database: " + err.Error())
	}

	view_db_dsn := os.Getenv("VIEW_DB_DSN")
	var view_db_conn *pgx.Conn = nil
	if view_db_dsn != "" {
		view_db_conn, err = connectToDB(view_db_dsn)
		if err != nil {
			log.Panicln("Failed to connect to view database: " + err.Error())
		}
	}
	return &RecordRepository{liveDb: live_db_conn, viewDb: view_db_conn}
}

func (repo *RecordRepository) GetNewRecords(sinceTime *time.Time) ([]models.Record, error) {
	var qry = ""
	if sinceTime == nil {
		qry = fmt.Sprintf(`SELECT id, revision, rel_module, "values", meta, rel_namespace, created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by FROM %s WHERE deleted_at is NULL;`, LIVE_DB_RECORD_TABLE)
	} else {
		qry = fmt.Sprintf(`SELECT id, revision, rel_module, "values", meta, rel_namespace, created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by FROM %s WHERE created_at > %s AND deleted_at is NULL;`, LIVE_DB_RECORD_TABLE, sinceTime.Format("2006-01-02 15:04:05-07"))
	}
	rows, err := repo.liveDb.Query(context.Background(), qry)
	if err != nil {
		log.Println("Failed to execute query: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var record models.Record
		if err := rows.Scan(&record.ID, &record.Revision, &record.ModuleID, &record.Values, &record.Meta, &record.NamespaceID, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt, &record.OwnedBy, &record.CreatedBy, &record.UpdatedBy, &record.DeletedBy); err != nil {
			log.Println("Failed to scan row: " + err.Error())
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (repo *RecordRepository) SyncNow() error {
	var latestCreatedAt time.Time

	query := fmt.Sprintf(`SELECT MAX(created_at) FROM %s`, VIEW_DB_INCIDENT_TABLE)

	err := repo.viewDb.QueryRow(context.Background(), query).Scan(&latestCreatedAt)
	if err != nil {
		log.Println("Failed to execute query: " + err.Error())
	}

	records, err := repo.GetNewRecords(&latestCreatedAt)
	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	var incidents []models.Incident
	for _, record := range records {
		incident := models.Incident{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}
		for _, value := range record.Values {
			switch value.Name {
			case "Channel":
				incident.Channel = value.Value
			case "Criticality":
				incident.Criticality = value.Value
			case "CallerName":
				incident.CallerName = value.Value
			case "LastCallDate":
				incident.LastCallDate = value.Value
			case "NationalID":
				incident.NationalID = value.Value
			case "MobileNumber":
				incident.MobileNumber = value.Value
			case "NotesOnCaller":
				incident.NotesOnCaller = value.Value
			case "IncidentReason":
				incident.IncidentReason = value.Value
			case "IncidentDescription":
				incident.IncidentDescription = value.Value
			case "Map":
				incident.Map = value.Value
			case "District":
				incident.District = value.Value
			case "Street":
				incident.Street = value.Value
			case "Location":
				incident.Location = value.Value
			case "IncidentNo":
				incident.IncidentNo = value.Value
			case "Status":
				incident.Status = value.Value
			}
		}

		incidents = append(incidents, incident)
	}
	log.Printf("Found %d new records to sync.\n", len(incidents))

	// TODO: Insert incidents into view database

	return nil
}
