package service

import (
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/models"
	"github.com/gokul-leadergroup/automax3.0_hermes/repository"
)

type recordService struct {
	repo *repository.RecordRepository
}

type RecordService interface {
	GetNewIncidents(sinceTime *time.Time) ([]models.Record, error)
}

func NewRecordService() RecordService {
	return &recordService{
		repo: repository.NewRecordRepository(),
	}
}

func (svc *recordService) GetNewIncidents(sinceTime *time.Time) ([]models.Record, error) {
	return svc.repo.GetNewRecords(sinceTime)
}
