package schemas

import (
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/pborman/uuid"
)

type MessageCreateSchema struct {
	ChatID   *uuid.UUID `json:"chat_id" validate:"omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	Username *string    `json:"username,omitempty" validate:"required_without=ChatID" example:"john-doe"`
	Text     *string    `json:"text" validate:"required_without=FileType" example:"I am not in danger skyler, I am the danger"`
	FileType *string    `json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
}

type MessageUpdateSchema struct {
	Text     *string `json:"text" validate:"required_without=FileType" example:"The Earth is the Lord's and the fullness thereof"`
	FileType *string `json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
}

type MessagesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []models.Message `json:"items"`
}

func (data MessagesResponseDataSchema) Init() MessagesResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type MessagesSchema struct {
	Chat     models.Chat                `json:"chat"`
	Messages MessagesResponseDataSchema `json:"messages"`
	Users    []models.UserDataSchema    `json:"users"`
}

func (data MessagesSchema) Init() MessagesSchema {
	// Set Initial Data
	// Set Chat
	chat := data.Chat.Init()
	data.Chat = chat
	// Set Users
	data.Users = chat.Users
	return data
}

type GroupChatInputSchema struct {
	Name              *string   `json:"name" validate:"omitempty,max=100" example:"Dopest Group"`
	Description       *string   `json:"description" validate:"omitempty,max=1000" example:"This is a group for bosses."`
	UsernamesToAdd    *[]string `json:"usernames_to_add" validate:"omitempty,min=1,max=99" example:"john-doe"`
	UsernamesToRemove *[]string `json:"usernames_to_remove" validate:"omitempty,min=1,max=99,usernames_to_update_validator" example:"john-doe"`
	FileType          *string   `json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
}

type GroupChatCreateSchema struct {
	Name           string   `json:"name" validate:"required,max=100" example:"Dopest Group"`
	Description    *string  `json:"description" validate:"omitempty,max=1000" example:"This is a group for bosses."`
	UsernamesToAdd []string `json:"usernames_to_add" validate:"required,min=1,max=99" example:"john-doe"`
	FileType       *string  `json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
}

// RESPONSE SCHEMAS
// CHATS
type ChatsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []models.Chat `json:"chats"`
}

func (data ChatsResponseDataSchema) Init() ChatsResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type ChatsResponseSchema struct {
	ResponseSchema
	Data ChatsResponseDataSchema `json:"data"`
}

type MessageCreateResponseSchema struct {
	ResponseSchema
	Data models.Message `json:"data"`
}

type ChatResponseSchema struct {
	ResponseSchema
	Data MessagesSchema `json:"data"`
}

type GroupChatInputResponseSchema struct {
	ResponseSchema
	Data models.Chat `json:"data"`
}
