package models

import (
	"fmt"

	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

type FeedAbstract struct {
	BaseModel
	AuthorID uuid.UUID `gorm:"not null"`
	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Text     string    `json:"text"`
	Slug     string    `gorm:"unique" json:"slug"`
}

func (obj *FeedAbstract) BeforeCreate(tx *gorm.DB) (err error) {
	// Create slug
	obj.Slug = fmt.Sprintf("%s-%s-%s", obj.Author.FirstName, obj.Author.LastName, obj.ID)
	return
}

type Post struct {
	FeedAbstract
	ImageID *uuid.UUID `gorm:"null"`
	Image   *File      `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL"`
}

type Comment struct {
	FeedAbstract
	PostID uuid.UUID `gorm:"not null"`
	Post   Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Reply struct {
	FeedAbstract
	CommentID uuid.UUID `gorm:"not null"`
	Comment   Comment   `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}

type Reaction struct {
	BaseModel
	UserID    uuid.UUID              `gorm:"not null;index:,unique,composite:user_id_post_id;index:,unique,composite:user_id_comment_id;index:,unique,composite:user_id_reply_id"`
	User      User                   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Rtype     choices.ReactionChoice `gorm:"varchar(50)"`
	PostID    uuid.UUID              `gorm:"null;index:,unique,composite:user_id_post_id"`
	Post      Post                   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentID uuid.UUID              `gorm:"null;index:,unique,composite:user_id_comment_id"`
	Comment   Comment                `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
	ReplyID   uuid.UUID              `gorm:"null;index:,unique,composite:user_id_reply_id"`
	Reply     Reply                  `gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE"`
}
