package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	Record struct {
		ID       uint64 `json:"recordID,string"`
		ModuleID uint64 `json:"moduleID,string"`
		Revision int    `json:"revision,omitempty"`

		Values RecordValueSet `json:"values,omitempty"`

		Meta map[string]any `json:"meta,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		OwnedBy         uint64     `json:"ownedBy,string"`
		CreatedAt       time.Time  `json:"createdAt,omitempty"`
		CreatedBy       uint64     `json:"createdBy,string"`
		UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy       uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt       *time.Time `json:"deletedAt,omitempty"`
		DeletedBy       uint64     `json:"deletedBy,string,omitempty"`
		RecordNumber    uint64     `json:"record_number"`
		RecordNumber007 uint64     `json:"record_number007"`
	}

	RecordValueSet []*RecordValue

	RecordValue struct {
		RecordID  uint64     `json:"-"`
		Name      string     `json:"name"`
		Value     string     `json:"value,omitempty"`
		Ref       uint64     `json:"-"`
		Place     uint       `json:"-"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		Updated  bool   `json:"-"`
		OldValue string `json:"-"`
	}
)

//
// ✅ Custom JSON Unmarshaler for RecordValueSet
//
func (rvs *RecordValueSet) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as array first
	var arr []*RecordValue
	if err := json.Unmarshal(data, &arr); err == nil {
		*rvs = arr
		return nil
	}

	// Try to unmarshal as object (key-value map)
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err == nil {
		for k, v := range obj {
			valStr := fmt.Sprintf("%v", v)
			*rvs = append(*rvs, &RecordValue{Name: k, Value: valStr})
		}
		return nil
	}

	return fmt.Errorf("invalid RecordValueSet JSON format: %s", string(data))
}
