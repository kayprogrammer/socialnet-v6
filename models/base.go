package models

import (
	"time"

	"github.com/pborman/uuid"
	"gorm.io/gorm"
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

func (f File) UpdateOrCreate (db *gorm.DB, id *uuid.UUID) File {
	if id == nil {
		db.Create(&f)
	} else {
		db.Model(File{BaseModel: BaseModel{ID: *id}}).Updates(&f)
		f.ID = *id
	}
	return f
}

type UserDataSchema struct {
	Name     string  `json:"name" example:"John Doe"`
	Username string  `json:"username" example:"john-doe"`
	Avatar   *string `json:"avatar" example:"https://img.url"`
}

func (user UserDataSchema) Init(userObj User) UserDataSchema {
	user.Name = userObj.FullName()
	user.Username = userObj.Username
	user.Avatar = userObj.GetAvatarUrl()
	return user
}
