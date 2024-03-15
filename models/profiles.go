package models

import "github.com/google/uuid"

type FriendStatusChoice string

const (
	PENDING  FriendStatusChoice = "PENDING"
	ACCEPTED FriendStatusChoice = "ACCEPTED"
)

type Friend struct {
	BaseModel
	RequesterID uuid.UUID          `gorm:"not null"`
	Requester   User               `gorm:"foreignKey:RequesterId;constraint:OnDelete:CASCADE"`
	RequesteeID uuid.UUID          `gorm:"not null"`
	Requestee   User               `gorm:"foreignKey:RequesteeId;constraint:OnDelete:CASCADE"`
	Status      FriendStatusChoice `gorm:"varchar(50)"`
}

type Notification struct {
	BaseModel
}
