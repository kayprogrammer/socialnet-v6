package managers

import (
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
type ChatManager struct {
}

func (obj ChatManager) GetUserChats(db *gorm.DB, userObj models.User) []models.Chat {
	chats := []models.Chat{}
	db.Where(models.Chat{OwnerID: userObj.ID}).Or(models.Chat{UserObjs: []models.User{userObj}}).Joins("OwnerObj").Joins("OwnerObj.AvatarObj").Joins("ImageObj").Joins("Messages")
	return chats
}

func (obj ChatManager) GetByID(db *gorm.DB, id uuid.UUID) models.Chat {
	chat := models.Chat{}
	db.Joins("Users").Take(&chat, models.Chat{BaseModel: models.BaseModel{ID: id}})
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
	db.Omit("ImageObj", "OwnerObj").Create(&chat)
	return chat
}

func (obj ChatManager) CreateGroup(db *gorm.DB, owner models.User, usersToAdd []models.User, data schemas.GroupChatCreateSchema) models.Chat {
	chat := models.Chat{
		OwnerID: owner.ID, 
		OwnerObj: owner, 
		Name: data.Name,
		Description: data.Description,
		Ctype: choices.CGROUP,
		UserObjs: usersToAdd,
	}

	fileType := data.FileType
	image := models.File{ResourceType: fileType}
	if fileType != nil {
		db.Create(&image)
		chat.ImageID = &image.ID
		chat.ImageObj = &image
	}
	db.Omit("ImageObj", "OwnerObj").Create(&chat)
	return chat
}

func (obj ChatManager) UsernamesToAddAndRemoveValidations(db *gorm.DB, chatObj models.Chat, chatUpdateQuery models.ChatUpdateOne, usernamesToAdd *[]string, usernamesToRemove *[]string) (models.ChatUpdateOne, *utils.ErrorResponse) {
	originalExistingUserIDs := []uuid.UUID{}
	for _, user := range chatObj.Edges.Users {
		originalExistingUserIDs = append(originalExistingUserIDs, user.ID)
	}
	expectedUserTotal := len(originalExistingUserIDs)
	usersToAdd := []models.User{}
	if usernamesToAdd != nil {
		usersToAdd = client.User.Query().
			Where(
				user.UsernameIn(*usernamesToAdd...),
				user.Or(
					user.Not(user.IDIn(originalExistingUserIDs...)),
					user.IDNEQ(chatObj.OwnerID),
				),
			).AllX(Ctx)
		expectedUserTotal += len(usersToAdd)
		chatUpdateQuery = chatUpdateQuery.AddUsers(usersToAdd...)
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
		usersToRemove = client.User.Query().
			Where(
				user.UsernameIn(*usernamesToRemove...),
				user.IDIn(originalExistingUserIDs...),
				user.IDNEQ(chatObj.OwnerID),
			).AllX(Ctx)
		expectedUserTotal -= len(usersToRemove)
		chatUpdateQuery = chatUpdateQuery.RemoveUsers(usersToRemove...)
	}
	if expectedUserTotal > 99 {
		data := map[string]string{
			"usernames_to_add": "99 users limit reached",
		}
		errData := utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)
		return nil, &errData
	}
	return chatUpdateQuery, nil
}

func (obj ChatManager) UpdateGroup(db *gorm.DB, chatObj models.Chat, data schemas.GroupChatInputSchema) (models.Chat, *utils.ErrorResponse) {
	chatUpdateQuery := chatObj.Update().
		SetNillableName(data.Name).
		SetNillableDescription(data.Description)

	// Handle users upload or remove
	var errData *utils.ErrorResponse
	chatUpdateQuery, errData = obj.UsernamesToAddAndRemoveValidations(client, chatObj, chatUpdateQuery, data.UsernamesToAdd, data.UsernamesToRemove)
	if errData != nil {
		return nil, errData
	}
	// Handle file upload
	var imageId *uuid.UUID
	image := chatObj.Edges.Image
	if data.FileType != nil {
		// Create or Update Image Object
		image = FileManager{}.UpdateOrCreate(client, image, *data.FileType)
		imageId = &image.ID
	}
	chatUpdateQuery = chatUpdateQuery.SetNillableImageID(imageId)
	updatedChat := chatUpdateQuery.SaveX(Ctx)

	// Set related data
	updatedChat.Edges.Users = chatObj.QueryUsers().WithAvatar().AllX(Ctx)
	updatedChat.Edges.Image = image

	return updatedChat, errData
}

func (obj ChatManager) GetSingleUserChat(db *gorm.DB, userObj models.User, id uuid.UUID) models.Chat {
	chat, _ := client.Chat.Query().
		Where(
			chat.IDEQ(id),
			chat.Or(
				chat.OwnerIDEQ(userObj.ID),
				chat.HasUsersWith(user.ID(userObj.ID)),
			),
		).
		Only(Ctx)
	return chat
}

func (obj ChatManager) GetSingleUserChatFullDetails(db *gorm.DB, userObj models.User, id uuid.UUID) models.Chat {
	chat, _ := client.Chat.Query().
		Where(
			chat.IDEQ(id),
			chat.Or(
				chat.OwnerIDEQ(userObj.ID),
				chat.HasUsersWith(user.ID(userObj.ID)),
			),
		).
		WithOwner(func(uq models.UserQuery) { uq.WithAvatar() }).
		WithImage().
		WithMessages(
			func(mq models.MessageQuery) {
				mq.WithSender(func(uq models.UserQuery) { uq.WithAvatar() }).WithFile().Order(ent.Desc(message.FieldCreatedAt))
			},
		).
		WithUsers(func(uq models.UserQuery) { uq.WithAvatar() }).
		Only(Ctx)
	return chat
}

func (obj ChatManager) GetUserGroup(db *gorm.DB, userObj models.User, id uuid.UUID, detailedOpts ...bool) models.Chat {
	chatQ := client.Chat.Query().
		Where(
			chat.CtypeEQ("GROUP"),
			chat.IDEQ(id),
			chat.OwnerIDEQ(userObj.ID),
		)
	if len(detailedOpts) > 0 {
		// Extra details
		chatQ = chatQ.
			WithOwner(func(uq models.UserQuery) { uq.WithAvatar() }).
			WithImage().
			WithUsers(func(uq models.UserQuery) { uq.WithAvatar() })
	}
	chatObj, _ := chatQ.Only(Ctx)
	return chatObj
}

func (obj ChatManager) GetMessagesCount(db *gorm.DB, chatID uuid.UUID) int {
	messagesCount := client.Message.Query().
		Where(
			message.ChatIDEQ(chatID),
		).CountX(Ctx)

	return messagesCount
}

func (obj ChatManager) DropData(db *gorm.DB) {
	client.Chat.Delete().ExecX(Ctx)
}

// ----------------------------------
// MESSAGE MANAGEMENT
// --------------------------------
type MessageManager struct {
}

func (obj MessageManager) Create(db *gorm.DB, sender models.User, chat models.Chat, text *string, fileType *string) models.Message {
	var fileID *uuid.UUID
	var file *ent.File
	if fileType != nil {
		file = FileManager{}.Create(client, *fileType)
		fileID = &file.ID
	}

	messageObj := client.Message.Create().
		SetChat(chat).
		SetSender(sender).
		SetNillableText(text).
		SetNillableFileID(fileID).
		SaveX(Ctx)

	// Set related values
	messageObj.Edges.Sender = sender
	if fileID != nil {
		messageObj.Edges.File = file
	}

	// Update Chat to intentionally update the updatedAt
	updatedChat := chat.Update().SaveX(Ctx)
	updatedChat.Edges.Owner = chat.Edges.Owner
	updatedChat.Edges.Image = chat.Edges.Image
	updatedChat.Edges.Users = chat.Edges.Users
	messageObj.Edges.Chat = updatedChat
	return messageObj
}

func (obj MessageManager) GetUserMessage(db *gorm.DB, userObj models.User, id uuid.UUID) models.Message {
	messageObj, _ := client.Message.Query().
		Where(
			message.IDEQ(id),
			message.SenderIDEQ(userObj.ID),
		).
		WithSender(func(uq models.UserQuery) { uq.WithAvatar() }).
		WithChat().
		WithFile().
		Only(Ctx)
	return messageObj
}

func (obj MessageManager) Update(db *gorm.DB, message models.Message, text *string, fileType *string) models.Message {
	var fileId *uuid.UUID
	file := message.Edges.File
	if fileType != nil {
		// Create or Update Image Object
		file = FileManager{}.UpdateOrCreate(client, file, *fileType)
		fileId = &file.ID
	}

	messageObj := message.Update().
		SetNillableText(text).
		SetNillableFileID(fileId).
		SaveX(Ctx)

	// Set related values
	messageObj.Edges.Sender = message.Edges.Sender
	if fileId != nil {
		messageObj.Edges.File = file
	}
	return messageObj
}

func (obj MessageManager) GetByID(db *gorm.DB, id uuid.UUID) models.Message {
	messageObj, _ := client.Message.Query().
		Where(
			message.IDEQ(id),
		).
		WithSender(func(uq models.UserQuery) { uq.WithAvatar() }).
		WithFile().
		Only(Ctx)
	return messageObj
}

func (obj MessageManager) DropData(db *gorm.DB) {
	client.Message.Delete().ExecX(Ctx)
}
