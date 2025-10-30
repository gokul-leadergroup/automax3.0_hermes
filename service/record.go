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
	GetNewRecords(sinceTime *time.Time) ([]models.Record, error)
	SyncNow() error
}

func NewRecordService() RecordService {
	return &recordService{
		repo: repository.NewRecordRepository(),
	}
}

func (svc *recordService) GetNewRecords(sinceTime *time.Time) ([]models.Record, error) {
	return svc.repo.GetNewRecords(sinceTime)
}

func (svc *recordService) SyncNow() error {
	return svc.repo.SyncNow()
}
