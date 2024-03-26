package models

import (
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
	ReactionsCount int        `json:"reactions_count" gorm:"-"`
}

type Post struct {
	FeedAbstract
	ImageID        *uuid.UUID             `gorm:"null" json:"-"`
	ImageObj       *File                  `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL" json:"-"`
	Image          *string                `gorm:"-" json:"image"`
	Comments       []Comment              `json:"-"`
	CommentsCount  int                    `json:"comments_count" gorm:"-"`
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
	PostID  uuid.UUID `json:"-" gorm:"not null"`
	PostObj Post      `json:"-" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Replies []Reply   `json:"-"`
	RepliesCount int   `json:"replies_count" gorm:"-" example:"50"`
}

func (c Comment) Init() Comment {
	c.ID = nil // Omit ID
	c.Author = UserDataSchema{}.Init(c.AuthorObj)
	c.RepliesCount = len(c.Replies)
	c.ReactionsCount = len(c.Reactions)
	return c
}

type Reply struct {
	FeedAbstract
	CommentID  uuid.UUID `json:"-" gorm:"not null"`
	CommentObj Comment   `json:"-" gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}

func (r Reply) Init() Reply {
	r.ID = nil // Omit ID
	r.Author = UserDataSchema{}.Init(r.AuthorObj)
	r.ReactionsCount = len(r.Reactions)
	return r
}
type Reaction struct {
	BaseModel
	UserID    uuid.UUID              `json:"-" gorm:"not null;index:,unique,composite:user_id_post_id;index:,unique,composite:user_id_comment_id;index:,unique,composite:user_id_reply_id"`
	UserObj   User                   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	User      UserDataSchema         `gorm:"-" json:"user"`
	Rtype     choices.ReactionChoice `gorm:"varchar(50)" json:"rtype" example:"LIKE"`
	PostID    *uuid.UUID              `json:"-" gorm:"null;index:,unique,composite:user_id_post_id"`
	Post      *Post                   `json:"-" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentID *uuid.UUID              `json:"-" gorm:"null;index:,unique,composite:user_id_comment_id"`
	Comment   *Comment                `json:"-" gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
	ReplyID   *uuid.UUID              `json:"-" gorm:"null;index:,unique,composite:user_id_reply_id"`
	Reply     *Reply                  `json:"-" gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE"`
}

func (r *Reaction) Init() {
	r.User = UserDataSchema{}.Init(r.UserObj)
}

func (r *Reaction) AfterFind(tx *gorm.DB) (err error) {
	r.Init()
	return
}

func (r *Reaction) AfterCreate(tx *gorm.DB) (err error) {
	r.Init()
	return
}
