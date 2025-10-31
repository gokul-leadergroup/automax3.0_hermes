package models

import (
	"time"
)

type (
	Classification struct {
		ID uint64 `json:"classificationID,string"`

		Name       string `json:"name"`
		ArabicName string `json:"arabicName"`

		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	}
)
