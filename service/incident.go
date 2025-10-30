package service

import "github.com/gokul-leadergroup/automax3.0_hermes/models"

type incidentService struct{}

type IncidentService interface{
	IncidentsFromRecords(records []models.Record) ([]models.Incident, error) 
}

func NewIncidentService() IncidentService {
	return &incidentService{}
}

func (svc *incidentService) IncidentsFromRecords(records []models.Record) ([]models.Incident, error) {
	var incidents []models.Incident
	for _, record := range records {
		incident := models.Incident{
			ID:                  record.ID,
			CreatedAt:           record.CreatedAt,
			UpdatedAt:           record.UpdatedAt,
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
	return incidents, nil
}
