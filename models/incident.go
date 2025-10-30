package models

import "time"

type Incident struct {
	ID 	  uint64  `json:"id"`
	Channel string  `json:"channel"`
	Criticality string  `json:"criticality"`
	CallerName string  `json:"caller_name"`
	LastCallDate string  `json:"last_call_date"`
	NationalID string  `json:"national_id"`
	MobileNumber string  `json:"mobile_number"`
	NotesOnCaller string  `json:"notes_on_caller"`
	IncidentReason string  `json:"incident_reason"`
	IncidentDescription string  `json:"incident_description"`
	Map string  `json:"map"`
	District string  `json:"district"`
	Street string  `json:"street"`
	AssignedTo []string  `json:"assigned_to"`
	TaskID uint64  `json:"task_id"`
	ClassificationID uint64  `json:"classification_id"`
	DepartmentID uint64  `json:"department_id"`
	PrimaryLocationID uint64  `json:"primary_location_id"`
	Location string  `json:"location"`
	IncidentNo string  `json:"incident_no"`
	Status string  `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at"`
}
