package service

import (
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/live_db_repository"
	"github.com/gokul-leadergroup/automax3.0_hermes/models"
)

type recordService struct {
	repo *live_db_repository.RecordRepository
}

type RecordService interface {
	GetNewRecords(sinceTime *time.Time) ([]models.Record, error)
	SyncNow() error
}

func NewRecordService() RecordService {
	return &recordService{
		repo: live_db_repository.NewRecordRepository(),
	}
}

func (svc *recordService) GetNewRecords(sinceTime *time.Time) ([]models.Record, error) {
	return svc.repo.GetNewRecords(sinceTime)
}

func (svc *recordService) SyncNow() error {
	return svc.repo.SyncNow()
}
