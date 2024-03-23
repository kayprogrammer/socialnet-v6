package models

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

type FeedAbstract struct {
	BaseModel
	AuthorID  uuid.UUID      `gorm:"not null" json:"-"`
	AuthorObj User           `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"-"`
	Author    UserDataSchema `gorm:"-" json:"author"`
	Text      string         `json:"text"`
	Slug      string         `gorm:"unique;not null;" json:"slug"`
	Reactions []Reaction     `json:"-"`
}

func (obj *FeedAbstract) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.Parse(uuid.New()) 
	obj.ID = id
	// Create slug
	obj.Slug = slug.Make(fmt.Sprintf("%s %s %s", obj.AuthorObj.FirstName, obj.AuthorObj.LastName, id))
	return
}

type Post struct {
	FeedAbstract
	ImageID        *uuid.UUID             `gorm:"null" json:"-"`
	ImageObj       *File                  `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL" json:"-"`
	Image          *string                `gorm:"-" json:"image"`
	Comments       []Comment              `json:"-"`
	Reactions      []Reaction             `json:"-"`
	CommentsCount  int                    `json:"comments_count"`
	ReactionsCount int                    `json:"reactions_count"`
	FileUploadData *utils.SignatureFormat `gorm:"-" json:"file_upload_data,omitempty"`
}

func (p Post) Init() Post {
	p.ID = nil // Omit ID
	p.Author = UserDataSchema{}.Init(p.AuthorObj)
	p.Image = p.GetImageUrl()
	p.CommentsCount = len(p.Comments)
	p.ReactionsCount = len(p.Reactions)
	return p
}

func (p Post) InitC(fileType *string) Post {
	// Updating response for when post is created
	p = p.Init()
	image := p.ImageObj
	if fileType != nil && image != nil { // Generate data when file is being uploaded
		fuData := utils.GenerateFileSignature(image.ID.String(), "posts")
		p.FileUploadData = &fuData
	}
	return p
}

func (p Post) GetImageUrl() *string {
	image := p.ImageObj
	if image != nil {
		url := utils.GenerateFileUrl(image.ID.String(), "posts", image.ResourceType)
		return &url
	}
	return nil
}

type Comment struct {
	FeedAbstract
	PostID  uuid.UUID `gorm:"not null"`
	PostObj Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Replies []Reply   `json:"-"`
}

type Reply struct {
	FeedAbstract
	CommentID  uuid.UUID `gorm:"not null"`
	CommentObj Comment   `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}

type Reaction struct {
	BaseModel
	UserID    uuid.UUID              `gorm:"not null;index:,unique,composite:user_id_post_id;index:,unique,composite:user_id_comment_id;index:,unique,composite:user_id_reply_id"`
	UserObj   User                   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Rtype     choices.ReactionChoice `gorm:"varchar(50)"`
	PostID    uuid.UUID              `gorm:"null;index:,unique,composite:user_id_post_id"`
	Post      Post                   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentID uuid.UUID              `gorm:"null;index:,unique,composite:user_id_comment_id"`
	Comment   Comment                `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
	ReplyID   uuid.UUID              `gorm:"null;index:,unique,composite:user_id_reply_id"`
	Reply     Reply                  `gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE"`
}
