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

const RECORD_TABLE = "compose_record_hab"

var recordRepo *RecordRepository

type RecordRepository struct {
	conn *pgx.Conn
}

func NewRecordRepository() *RecordRepository {
	if recordRepo != nil {
		return recordRepo
	}

	view_db_dsn := os.Getenv("VIEW_DB_DSN")
	if view_db_dsn != "" {
		view_db_conn, err := pgx.Connect(context.Background(), view_db_dsn)
		if err != nil {
			log.Panicln("Failed to connect to view database: " + err.Error())
		}
		recordRepo = &RecordRepository{conn: view_db_conn}
		return recordRepo
	} else {
		log.Println("ENV variable not set for VIEW_DB_DSN")
		return nil
	}
}

func (repo *RecordRepository) LatestRecordCreatedAt() (*time.Time, error) {
	var latestCreatedAt *time.Time
	query := fmt.Sprintf(`SELECT MAX(created_at) FROM %s`, RECORD_TABLE)
	err := repo.conn.QueryRow(context.Background(), query).Scan(&latestCreatedAt)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	return latestCreatedAt, nil
}

func (repo *RecordRepository) BulkInsert(records []models.Record, tx pgx.Tx) error {
	if len(records) == 0 {
		log.Println("No new records to insert.")
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
	`, RECORD_TABLE)

	failed := []models.Incident{}
	for _, incident := range incidents {
		_, err := tx.Exec(context.Background(), insertLiveDbQry,
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
		log.Printf("Failed to insert %d incidents.\n", len(failed))
	} else {
		log.Println("All incidents insertec successfully.")
	}

	return nil
}
