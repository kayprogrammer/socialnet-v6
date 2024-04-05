package managers

import (
	"log"

	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

// ----------------------------------
// CHAT MANAGEMENT
// --------------------------------
func ChatOwnerImageScope(db *gorm.DB) *gorm.DB {
	return db.Joins("OwnerObj").Joins("OwnerObj.AvatarObj").Joins("ImageObj")
}

func MessageSenderFileScope(db *gorm.DB) *gorm.DB {
	return db.Joins("SenderObj").Joins("SenderObj.AvatarObj").Joins("FileObj")
}

func ChatPreloadMessagesScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Messages", func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(MessageSenderFileScope).Order("messages.created_at DESC")
	})
}

type ChatManager struct {
}

func (obj ChatManager) GetUserChats(db *gorm.DB, user models.User) []models.Chat {
	chats := []models.Chat{}
	db.Model(&models.Chat{}).
		Where(models.Chat{OwnerID: user.ID}).
		Or("chats.id IN (?)", db.Table("chat_users").Select("chat_id").Where("user_id = ?", user.ID)).
		Scopes(ChatOwnerImageScope, ChatPreloadMessagesScope).
		Find(&chats)
	return chats
}

func (obj ChatManager) GetByID(db *gorm.DB, id uuid.UUID) models.Chat {
	chat := models.Chat{}
	db.Preload("UserObjs").Take(&chat, models.Chat{BaseModel: models.BaseModel{ID: id}})
	return chat
}

func (obj ChatManager) UserIsMember(chat models.Chat, targetUser models.User) bool {
	for _, user := range chat.UserObjs {
		if user.ID.String() == targetUser.ID.String() {
			return true
		}
	}
	return false
}

func (obj ChatManager) GetDMChat(db *gorm.DB, user models.User, recipientUser models.User) models.Chat {
	chat := models.Chat{Ctype: choices.CDM}
	db.Where(models.Chat{OwnerID: user.ID, UserObjs: []models.User{recipientUser}}).Or(models.Chat{OwnerID: recipientUser.ID, UserObjs: []models.User{user}}).Take(&chat, chat)
	return chat
}

func (obj ChatManager) Create(db *gorm.DB, owner models.User, ctype choices.ChatTypeChoice, recipientsOpts ...[]models.User) models.Chat {
	chat := models.Chat{Ctype: ctype, OwnerID: owner.ID, OwnerObj: owner}
	if len(recipientsOpts) > 0 {
		chat.UserObjs = recipientsOpts[0]
	}
	db.Create(&chat)
	return chat
}

func (obj ChatManager) CreateGroup(db *gorm.DB, owner models.User, usersToAdd []models.User, data schemas.GroupChatCreateSchema) models.Chat {
	chat := models.Chat{
		OwnerID:     owner.ID,
		OwnerObj:    owner,
		Name:        &data.Name,
		Description: data.Description,
		Ctype:       choices.CGROUP,
		UserObjs:    usersToAdd,
	}

	fileType := data.FileType
	if fileType != nil {
		var fileType string = *data.FileType
		image := models.File{ResourceType: fileType}
		db.Create(&image)
		chat.ImageID = &image.ID
		chat.ImageObj = &image
	}
	db.Omit("UserObjs.*").Create(&chat)
	return chat
}

func (obj ChatManager) GetByUsernames(db *gorm.DB, usernames []string, excludeOpts ...uuid.UUID) []models.User {
	users := []models.User{}
	usersQ := db.Where("username IN ?", usernames)
	if len(excludeOpts) > 0 {
		usersQ = usersQ.Not("id = ?", excludeOpts[0])
	}
	usersQ.Find(&users)
	return users
}

func (obj ChatManager) UsernamesToAddAndRemoveValidations(db *gorm.DB, chat *models.Chat, usernamesToAdd *[]string, usernamesToRemove *[]string) (*models.Chat, *utils.ErrorResponse) {
	originalExistingUserIDs := []uuid.UUID{}
	for _, user := range chat.UserObjs {
		originalExistingUserIDs = append(originalExistingUserIDs, user.ID)
	}
	expectedUserTotal := len(originalExistingUserIDs)
	usersToAdd := []models.User{}
	if usernamesToAdd != nil {
		db.Where("username IN ?", usernamesToAdd).Where(
			db.Not("id IN ?", originalExistingUserIDs).Or(models.User{BaseModel: models.BaseModel{ID: chat.OwnerID}}),
		).Find(&usersToAdd)
		expectedUserTotal += len(usersToAdd)
	}
	usersToRemove := []models.User{}
	if usernamesToRemove != nil {
		if len(originalExistingUserIDs) < 1 {
			data := map[string]string{
				"usernames_to_remove": "No users to remove",
			}
			errData := utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)
			return nil, &errData
		}
		db.Where("username IN ?", usernamesToRemove).Not(models.User{BaseModel: models.BaseModel{ID: chat.OwnerID}}).Find(&usernamesToRemove, originalExistingUserIDs)
		expectedUserTotal -= len(usersToRemove)
	}
	if expectedUserTotal > 99 {
		data := map[string]string{
			"usernames_to_add": "99 users limit reached",
		}
		errData := utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)
		return nil, &errData
	}
	db.Model(&chat).Association("UserObjs").Append(&usersToAdd)
	db.Model(&chat).Association("UserObjs").Delete(&usernamesToRemove)
	return chat, nil
}

func (obj ChatManager) UpdateGroup(db *gorm.DB, chat *models.Chat, data schemas.GroupChatInputSchema) (*models.Chat, *utils.ErrorResponse) {
	if data.Name != nil {
		chat.Name = data.Name
	}
	if data.Description != nil {
		chat.Description = data.Description
	}

	// Handle users upload or remove
	var errData *utils.ErrorResponse
	chat, errData = obj.UsernamesToAddAndRemoveValidations(db, chat, data.UsernamesToAdd, data.UsernamesToRemove)
	if errData != nil {
		return nil, errData
	}
	// Handle file upload
	if data.FileType != nil {
		// Create or Update Image Object
		image := models.File{ResourceType: *data.FileType}.UpdateOrCreate(db, chat.ImageID)
		chat.ImageID = &image.ID
		chat.ImageObj = &image
	}
	db.Save(&chat)
	return chat, errData
}

func (obj ChatManager) GetSingleUserChat(db *gorm.DB, user models.User, id uuid.UUID) models.Chat {
	chat := models.Chat{}
	db.Model(&models.Chat{}).Where(models.Chat{OwnerID: user.ID}).
		Or("chats.id IN (?)", db.Table("chat_users").Select("chat_id").Where("user_id = ?", user.ID)).
		Take(&chat, id)
	return chat
}

func (obj ChatManager) GetSingleUserChatFullDetails(db *gorm.DB, user models.User, id uuid.UUID) models.Chat {
	chat := models.Chat{} // Wahala wa o
	db.Model(&models.Chat{}).Where(models.Chat{OwnerID: user.ID}).
		Or("chats.id IN (?)", db.Table("chat_users").Select("chat_id").Where("user_id = ?", user.ID)).
		Scopes(ChatOwnerImageScope, ChatPreloadMessagesScope).
		Preload("UserObjs").
		Take(&chat, id)
	return chat
}

func (obj ChatManager) GetUserGroup(db *gorm.DB, user models.User, id uuid.UUID, detailedOpts ...bool) models.Chat {
	chat := models.Chat{Ctype: choices.CGROUP, OwnerID: user.ID}
	q := db
	if len(detailedOpts) > 0 {
		q = q.Scopes(ChatOwnerImageScope).Preload("UserObjs")
	}
	q.Where(&chat).Take(&chat, id)
	log.Println(chat)
	return chat
}

func (obj ChatManager) GetMessagesCount(db *gorm.DB, chatID uuid.UUID) int64 {
	var messagesCount int64
	db.Model(&models.Message{ChatID: chatID}).Count(&messagesCount)
	return messagesCount
}

func (obj ChatManager) DropData(db *gorm.DB) {
	db.Delete(models.Chat{})
}

// ----------------------------------
// MESSAGE MANAGEMENT
// --------------------------------

func MessageSenderScope(db *gorm.DB) *gorm.DB {
	return db.Joins("SenderObj").Joins("SenderObj.AvatarObj").Joins("ChatObj").Joins("FileObj")
}

type MessageManager struct {
}

func (obj MessageManager) Create(db *gorm.DB, sender models.User, chat models.Chat, text *string, fileType *string) models.Message {
	message := models.Message{SenderID: sender.ID, SenderObj: sender, ChatID: chat.ID, ChatObj: chat, Text: text}
	if fileType != nil {
		file := models.File{ResourceType: *fileType}
		db.Create(&file)
		message.FileID = &file.ID
		message.FileObj = &file
	}
	db.Create(&message)
	return message
}

func (obj MessageManager) GetUserMessage(db *gorm.DB, user models.User, id uuid.UUID) models.Message {
	message := models.Message{SenderID: user.ID}
	db.Scopes(MessageSenderScope).Take(&message, models.Message{BaseModel: models.BaseModel{ID: id}})
	return message
}

func (obj MessageManager) Update(db *gorm.DB, message models.Message, text *string, fileType *string) models.Message {
	if fileType != nil {
		// Create or Update Image Object
		file := models.File{ResourceType: *fileType}.UpdateOrCreate(db, message.FileID)
		message.FileID = &file.ID
		message.FileObj = &file
	}
	if text != nil {
		message.Text = text
	}
	db.Save(&message)
	return message
}

func (obj MessageManager) GetByID(db *gorm.DB, id uuid.UUID) models.Message {
	message := models.Message{}
	db.Scopes(MessageSenderScope).Take(&message, models.Message{BaseModel: models.BaseModel{ID: id}})
	return message
}

// func (obj MessageManager) DropData(db *gorm.DB) {
// 	client.Message.Delete().ExecX(Ctx)
// }
