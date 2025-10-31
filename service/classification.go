package service

import (
	"time"

	"github.com/gokul-leadergroup/automax3.0_hermes/models"
	"github.com/gokul-leadergroup/automax3.0_hermes/repository"
)

type classificationService struct {
	repo *repository.ClassificationRepository
}

type ClassificationService interface {
	GetNewClassifications(sinceTime *time.Time) ([]models.Classification, error)
	SyncNow() error
}

func NewClassificationService() ClassificationService {
	return &classificationService{
		repo: repository.NewClassificationRepository(),
	}
}

func (svc *classificationService) GetNewClassifications(sinceTime *time.Time) ([]models.Classification, error) {
	return svc.repo.GetNewClassifications(sinceTime)
}

func (svc *classificationService) SyncNow() error {
	return svc.repo.SyncNow()
}
