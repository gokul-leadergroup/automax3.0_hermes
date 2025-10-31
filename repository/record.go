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
const VIEW_DB_INCIDENT_TABLE = "compose_record_hab"

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
		qry = fmt.Sprintf(`
		SELECT id, revision, rel_module, "values", meta, rel_namespace, 
		       created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by 
		FROM %s 
		WHERE deleted_at IS NULL;
	`, LIVE_DB_RECORD_TABLE)
	} else {
		qry = fmt.Sprintf(`
		SELECT id, revision, rel_module, "values", meta, rel_namespace, 
		       created_at, updated_at, deleted_at, owned_by, created_by, updated_by, deleted_by 
		FROM %s 
		WHERE created_at > '%s' AND deleted_at IS NULL;
	`, LIVE_DB_RECORD_TABLE, sinceTime.Format("2006-01-02 15:04:05-07"))
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
	var latestCreatedAt *time.Time // use a pointer

	query := fmt.Sprintf(`SELECT MAX(created_at) FROM %s`, VIEW_DB_INCIDENT_TABLE)
	err := repo.viewDb.QueryRow(context.Background(), query).Scan(&latestCreatedAt)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return err
	}

	if latestCreatedAt == nil {
		log.Println("Table is empty — no records found.")
	} else {
		log.Println("Latest created_at:", latestCreatedAt)
	}

	var records []models.Record
	records, err = repo.GetNewRecords(latestCreatedAt)

	if err != nil {
		log.Println("Failed to get new records: " + err.Error())
		return err
	}

	if len(records) == 0 {
		log.Println("No new records to sync.")
		return nil
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

	insertLiveDbQry := fmt.Sprintf(`
		INSERT INTO %s (
			id, module_id, namespace_id, created_at, created_by, updated_at, updated_by, status, channel, district,
			location, department, assigned_to, caller_name, criticality, mobile_number, classification, primary_location,
			incident_description, incident_no, national_id, street
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22
		);
	`, VIEW_DB_INCIDENT_TABLE)

	failed := []models.Incident{}
	for _, incident := range incidents {
		_, err := repo.viewDb.Exec(context.Background(), insertLiveDbQry,
			incident.ID,
			0,
			0,
			incident.CreatedAt,
			0,
			incident.UpdatedAt,
			0,
			incident.Status,
			incident.Channel,
			incident.District,
			incident.Location,
			incident.DepartmentID,
			incident.AssignedTo,
			incident.CallerName,
			incident.Criticality,
			incident.MobileNumber,
			incident.ClassificationID,
			incident.PrimaryLocationID,
			incident.IncidentDescription,
			incident.IncidentNo,
			incident.NationalID,
			incident.Street,
		)
		if err != nil {
			log.Println("Failed to insert incident: " + err.Error())
			failed = append(failed, incident)
		}
	}

	if len(failed) > 0 {
		log.Printf("Failed to sync %d incidents.\n", len(failed))
	} else {
		log.Println("All incidents synced successfully.")
	}

	return nil
}
