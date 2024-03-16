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
	Requester   User               `gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	RequesteeID uuid.UUID          `gorm:"not null;check:requester_id <> requestee_id"`
	Requestee   User               `gorm:"foreignKey:RequesteeID;constraint:OnDelete:CASCADE"`
	Status      FriendStatusChoice `gorm:"varchar(50)"`
}

type NotificationChoice string

const (
	REACTION NotificationChoice = "REACTION"
	COMMENT  NotificationChoice = "COMMENT"
	REPLY    NotificationChoice = "REPLY"
	ADMIN    NotificationChoice = "ADMIN"
)

type Notification struct {
	BaseModel
	SenderID  uuid.UUID          `gorm:"not null"`
	Sender    User               `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Ntype     NotificationChoice `gorm:"varchar(50);not null"`
	PostID    uuid.UUID          `gorm:"null"`
	Post      Post               `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentID uuid.UUID          `gorm:"null"`
	Comment   Comment            `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
	ReplyID   uuid.UUID          `gorm:"null"`
	Reply     Reply              `gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE"`
}
