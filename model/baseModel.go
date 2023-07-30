package model

import (
	"time"

	"gorm.io/gorm"
)

type CreatedBase struct {
	CreatedBy string    `gorm:"size:100;null;" json:"created_by"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

type DeletedBase struct {
	DeletedBy string         `gorm:"size:100;null;" json:"deleted_by"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type ArchievedBase struct {
	ArchievedAt time.Time `gorm:"default:current_timestamp;null;" json:"archieved_at"`
}
