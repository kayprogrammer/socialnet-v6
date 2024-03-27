package models

import (
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
)

type Chat struct {
	BaseModel
	OwnerID     uuid.UUID              `json:"-"`
	OwnerObj    User                   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Owner       UserDataSchema         `gorm:"-" json:"owner"`
	Name        *string                `gorm:"varchar(50)" json:"name" example:"My Group"`
	Ctype       choices.ChatTypeChoice `gorm:"varchar(50);check:(ctype = 'GROUP' AND name IS NOT NULL) OR (ctype = 'DM')" json:"ctype" example:"DM"`
	Description *string                `gorm:"varchar(200);check:(ctype = 'DM' AND name IS NULL AND description IS NULL AND image_id IS NULL) OR (ctype = 'GROUP')" json:"description" example:"My Group Description"`
	ImageID     *uuid.UUID             `json:"-"`
	ImageObj    *File                  `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL" json:"-"`
	Image       *string                `gorm:"-" json:"image"`
	Messages	[]Message              `gorm:"-"`
}

func (c Chat) GetImageUrl() *string {
	image := c.ImageObj
	if image != nil {
		url := utils.GenerateFileUrl(image.ID.String(), "chats", image.ResourceType)
		return &url
	}
	return nil
}

type Message struct {
	BaseModel
	SenderID  uuid.UUID      `json:"-"`
	SenderObj User           `json:"-" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Sender    UserDataSchema `gorm:"-" json:"sender"`
	ChatID    uuid.UUID      `json:"-"`
	ChatObj   Chat           `json:"-" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	Text      *string        `gorm:"varchar(1000000)" json:"text" example:"My Message"`
	FileID    *uuid.UUID     `json:"-"`
	FileObj   *File          `gorm:"foreignKey:FileID;constraint:OnDelete:SET NULL" json:"-"`
	File      *string        `gorm:"-" json:"image"`
}
