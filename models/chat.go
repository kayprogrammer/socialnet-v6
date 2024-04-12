package models

import (
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

type LatestMessageSchema struct {
	Sender UserDataSchema `json:"sender"`
	Text   *string        `json:"text"`
	File   *string        `json:"file"`
}

type Chat struct {
	BaseModel
	OwnerID     uuid.UUID              `json:"-"`
	OwnerObj    User                   `json:"-" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;<-:false"`
	Owner       UserDataSchema         `gorm:"-" json:"owner"`
	Name        *string                `gorm:"varchar(50)" json:"name" example:"My Group"`
	Ctype       choices.ChatTypeChoice `gorm:"varchar(50);check:(ctype = 'GROUP' AND name IS NOT NULL) OR (ctype = 'DM')" json:"ctype" example:"DM"`
	Description *string                `gorm:"varchar(200);check:(ctype = 'DM' AND name IS NULL AND description IS NULL AND image_id IS NULL) OR (ctype = 'GROUP')" json:"description" example:"A nice group for tech enthusiasts"`
	ImageID     *uuid.UUID             `json:"-"`
	ImageObj    *File                  `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL;<-:false" json:"-"`
	Image       *string                `gorm:"-" json:"image" example:"https://img.url"`
	UserObjs    []User                 `json:"-" gorm:"many2many:chat_users;"`
	Messages    []Message              `json:"-"`

	LatestMessage *LatestMessageSchema `gorm:"-" json:"latest_message"`
	Users		[]UserDataSchema		`gorm:"-" json:"users,omitempty" swaggerIgnore:"true"` // omitempty later to show for groups
	FileUploadData *utils.SignatureFormat `gorm:"-" json:"file_upload_data,omitempty"`
}

func (c *Chat) BeforeDelete (tx *gorm.DB) (err error) {
	tx.Model(&c).Association("UserObjs").Clear()
	return
}

func (c Chat) Init() Chat {
	// Set Owner Details.
	c.Owner = c.Owner.Init(c.OwnerObj)

	// Set ImageUrl
	c.Image = c.GetImageUrl()

	// Set Latest Message
	latestMessages := c.Messages
	if len(latestMessages) > 0 {
		latestMessage := latestMessages[0]
		file := latestMessage.FileObj
		var fileUrl *string
		if file != nil {
			url := utils.GenerateFileUrl(file.ID.String(), "messages", file.ResourceType)
			fileUrl = &url
		}
		lm := LatestMessageSchema{
			Text: latestMessage.Text,
			File: fileUrl,
		}
		lm.Sender = lm.Sender.Init(latestMessage.SenderObj)
		c.LatestMessage = &lm
	}
	return c
}

func (c Chat) InitG() Chat {
	// Init Group Chat
	c = c.Init()
	// Set Users Details for groups.
	users := []UserDataSchema{}
	for _, user := range c.UserObjs {
		userData := UserDataSchema{}.Init(user)
		users = append(users, userData)
	}
	c.Users = users
	return c
}


func (c Chat) GetImageUrl() *string {
	image := c.ImageObj
	if image != nil {
		url := utils.GenerateFileUrl(image.ID.String(), "groups", image.ResourceType)
		return &url
	}
	return nil
}

func (c Chat) InitC(fileType *string) Chat {
	c = c.InitG()
	// When chat is created
	file := c.ImageObj
	if fileType != nil && file != nil { // Generate data when file is being uploaded
		fuData := utils.GenerateFileSignature(file.ID.String(), "groups")
		c.FileUploadData = &fuData
	}
	return c
}

type Message struct {
	BaseModel
	SenderID  uuid.UUID      `json:"-"`
	SenderObj User           `json:"-" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;<-:false;"`
	Sender    UserDataSchema `gorm:"-" json:"sender"`
	ChatID    uuid.UUID      `json:"chat_id"`
	ChatObj   Chat           `json:"-" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE;<-:false"`
	Text      *string        `gorm:"varchar(1000000)" json:"text" example:"Jesus is King"`
	FileID    *uuid.UUID     `json:"-"`
	FileObj   *File          `gorm:"foreignKey:FileID;constraint:OnDelete:SET NULL;<-:false" json:"-"`
	File      *string        `gorm:"-" json:"file" example:"https://img.url"`
	FileUploadData *utils.SignatureFormat `gorm:"-" json:"file_upload_data,omitempty"`
}

func (m *Message) AfterCreate(tx *gorm.DB) (err error) {
	// Update Chat to intentionally update the updatedAt
	tx.Save(&m.ChatObj)
	return
}

func (m Message) Init() Message {
	// Set Author Details.
	m.Sender = m.Sender.Init(m.SenderObj)

	// Set FileUrl
	file := m.FileObj
	if file != nil {
		url := utils.GenerateFileUrl(file.ID.String(), "messages", file.ResourceType)
		m.File = &url
	}
	return m
}

func (m Message) InitC(fileType *string) Message {
	m = m.Init()
	// When message is created
	file := m.FileObj
	if fileType != nil && file != nil { // Generate data when file is being uploaded
		fuData := utils.GenerateFileSignature(file.ID.String(), "messages")
		m.FileUploadData = &fuData
	}
	return m
}