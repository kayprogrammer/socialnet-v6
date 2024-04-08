package models

import (
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

type Friend struct {
	BaseModel
	RequesterID uuid.UUID                  `gorm:"not null"`
	Requester   User                       `gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	RequesteeID uuid.UUID                  `gorm:"not null;check:requester_id <> requestee_id"`
	Requestee   User                       `gorm:"foreignKey:RequesteeID;constraint:OnDelete:CASCADE"`
	Status      choices.FriendStatusChoice `gorm:"varchar(50)"`
}

type Notification struct {
	BaseModel
	SenderID  *uuid.UUID                 `gorm:"null" json:"-"`
	SenderObj *User                      `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;<-:false" json:"-"`
	Sender    *UserDataSchema            `gorm:"-" json:"sender"`
	Receivers []User                     `gorm:"many2many:notification_receivers;" json:"-"`
	Ntype     choices.NotificationChoice `json:"ntype" gorm:"varchar(50);not null"`
	Text      *string                    `gorm:"varchar(10000);null;" json:"-"`
	PostID    *uuid.UUID                 `json:"-" gorm:"null"`
	Post      *Post                      `json:"-" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;<-:false"`
	CommentID *uuid.UUID                 `json:"-" gorm:"null"`
	Comment   *Comment                   `json:"-" gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE;<-:false"`
	ReplyID   *uuid.UUID                 `json:"-" gorm:"null"`
	Reply     *Reply                     `json:"-" gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE;<-:false"`
	ReadBy    []User                     `json:"-" gorm:"many2many:notification_read_by;"`

	// Other schema display
	PostSlug    *string `gorm:"-" json:"post_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	CommentSlug *string `gorm:"-" json:"comment_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	ReplySlug   *string `gorm:"-" json:"reply_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	Message     string  `gorm:"-" json:"message" example:"John Doe reacted to your post"`
	IsRead      bool    `gorm:"-" json:"is_read" example:"true"`
}

func (n Notification) BeforeDelete (tx *gorm.DB) (err error) {
	tx.Model(n).Association("Receivers").Clear()
	return
}

func (n Notification) Init(currentUserID uuid.UUID) Notification {
	// Set Related Data.
	sender := n.SenderObj
	if sender != nil {
		senderData := UserDataSchema{}.Init(*sender)
		n.Sender = &senderData
	}

	// Set Target slug
	n = n.SetTargetSlug()
	// Set Notification message
	text := n.Text
	if text == nil {
		notificationMsg := n.GetMessage()
		text = &notificationMsg
	}
	n.Message = *text

	// Set IsRead
	if currentUserID != nil {
		readBy := n.ReadBy
		for _, user := range readBy {
			if user.ID.String() == currentUserID.String() {
				n.IsRead = true
				break
			}
		}
	}
	return n
}

func (n Notification) SetTargetSlug() Notification {
	post := n.Post
	comment := n.Comment
	reply := n.Reply
	if post != nil {
		n.PostSlug = &post.Slug
	} else if comment != nil {
		n.CommentSlug = &comment.Slug
	} else if reply != nil {
		n.ReplySlug = &reply.Slug
	}
	return n

}

func (n Notification) GetMessage() string {
	ntype := n.Ntype
	sender := n.Sender.Name
	message := sender + " reacted to your post"
	if ntype == "REACTION" {
		if n.CommentSlug != nil {
			message = sender + " reacted to your comment"
		} else if n.ReplySlug != nil {
			message = sender + " reacted to your reply"
		}
	} else if ntype == "COMMENT" {
		message = sender + " commented on your post"
	} else if ntype == "REPLY" {
		message = sender + " replied your comment"
	}
	return message
}
