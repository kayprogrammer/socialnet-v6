package models

import (
	"time"

	"github.com/pborman/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primarykey;not null;default:uuid_generate_v4()" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type File struct {
	BaseModel
	ResourceType string `json:"resource_type" gorm:"not null"`
}
