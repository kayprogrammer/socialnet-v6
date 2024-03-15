package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedAbstract struct {
	BaseModel
	AuthorID uuid.UUID
	Author   User   `gorm:"foreignKey:AuthorId;constraint:OnDelete:CASCADE"`
	Text     string `json:"text"`
	Slug     string `gorm:"unique" json:"slug"`
}

func (obj *FeedAbstract) BeforeCreate(tx *gorm.DB) (err error) {
	// Create slug
	obj.Slug = fmt.Sprintf("%v-%v-%v", obj.Author.FirstName, obj.Author.LastName, obj.ID)
	return
}

type Post struct {
	FeedAbstract
	ImageID *uuid.UUID `gorm:"null"`
	Image   *File      `gorm:"foreignKey:ImageId;constraint:OnDelete:SET NULL"`
}

type Comment struct {
	FeedAbstract
	PostID uuid.UUID `gorm:"not null"`
	Post   Post      `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
}

type Reply struct {
	FeedAbstract
	CommentID uuid.UUID `gorm:"not null"`
	Comment   Comment   `gorm:"foreignKey:CommentId;constraint:OnDelete:CASCADE"`
}

type ReactionChoice string

const (
	LIKE  ReactionChoice = "LIKE"
	LOVE  ReactionChoice = "LOVE"
	HAHA  ReactionChoice = "HAHA"
	WOW   ReactionChoice = "WOW"
	SAD   ReactionChoice = "SAD"
	ANGRY ReactionChoice = "ANGRY"
)

type Reaction struct {
	BaseModel
	UserID    uuid.UUID      `gorm:"not null;index:,unique,composite:user_id_post_id;index:,unique,composite:user_id_comment_id;index:,unique,composite:user_id_reply_id"`
	User      User           `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Rtype     ReactionChoice `gorm:"varchar(50)"`
	PostID    uuid.UUID      `gorm:"null;index:,unique,composite:user_id_post_id"`
	Post      Post           `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
	CommentID uuid.UUID      `gorm:"null;index:,unique,composite:user_id_comment_id"`
	Comment   Comment        `gorm:"foreignKey:CommentId;constraint:OnDelete:CASCADE"`
	ReplyID   uuid.UUID      `gorm:"null;index:,unique,composite:user_id_reply_id"`
	Reply     Reply          `gorm:"foreignKey:ReplyId;constraint:OnDelete:CASCADE"`
}
