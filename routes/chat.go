package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v4/ent"
	"github.com/kayprogrammer/socialnet-v4/managers"
	"github.com/kayprogrammer/socialnet-v4/schemas"
	"github.com/kayprogrammer/socialnet-v4/utils"
)

var (
	chatManager    = managers.ChatManager{}
	messageManager = managers.MessageManager{}
)

// @Summary Retrieve User Chats
// @Description `This endpoint retrieves a paginated list of the current user chats`
// @Tags Chat
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.ChatsResponseSchema
// @Router /chats [get]
// @Security BearerAuth
func (endpoint Endpoint) RetrieveUserChats(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)
	chats := chatManager.GetUserChats(db, user)

	// Paginate, Convert type and return chats
	paginatedData, paginatedChats, err := PaginateQueryset(chats, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	convertedChats := utils.ConvertStructData(paginatedChats, []schemas.ChatSchema{}).(*[]schemas.ChatSchema)
	response := schemas.ChatsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Chats fetched"}.Init(),
		Data: schemas.ChatsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       *convertedChats,
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Send a message
// @Description `This endpoint sends a message`
// @Description
// @Description `You must either send a text or a file or both.`
// @Description
// @Description `If there's no chat_id, then its a new chat and you must set username and leave chat_id`
// @Description
// @Description `If chat_id is available, then ignore username and set the correct chat_id`
// @Description
// @Description `The file_upload_data in the response is what is used for uploading the file to cloudinary from client`
// @Tags Chat
// @Param message body schemas.MessageCreateSchema true "Message object"
// @Success 201 {object} schemas.MessageCreateResponseSchema
// @Router /chats [post]
// @Security BearerAuth
func (endpoint Endpoint) SendMessage(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)

	messageData := schemas.MessageCreateSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &messageData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(messageData); err != nil {
		return c.Status(422).JSON(err)
	}

	chatID := messageData.ChatID
	username := messageData.Username

	var chat *ent.Chat
	if chatID == nil {
		// Create a new chat dm with current user and recipient user
		recipientUser := userManager.GetByUsername(db, *username)
		if recipientUser == nil {
			data := map[string]string{
				"username": "No user with that username",
			}
			return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid entry", data))
		}
		chat = chatManager.GetDMChat(db, user, recipientUser)
		// Check if a chat already exists between both users
		if chat != nil {
			data := map[string]string{
				"username": "A chat already exist between you and the recipient",
			}
			return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid entry", data))
		}
		chat = chatManager.Create(db, user, "DM", []*ent.User{recipientUser})
	} else {
		// Get the chat with chat id and check if the current user is the owner or the recipient
		chat = chatManager.GetSingleUserChat(db, user, *chatID)
		if chat == nil {
			return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no chat with that ID"))
		}
	}

	//Create Message
	message := messageManager.Create(db, user, chat, messageData.Text, messageData.FileType)

	// Convert type and return Message
	convertedMessage := utils.ConvertStructData(message, schemas.MessageCreateResponseDataSchema{}).(*schemas.MessageCreateResponseDataSchema)
	response := schemas.MessageCreateResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Message sent"}.Init(),
		Data:           convertedMessage.Init(messageData.FileType),
	}
	return c.Status(201).JSON(response)
}

// @Summary Retrieve messages from a Chat
// @Description `This endpoint retrieves all messages in a chat`
// @Tags Chat
// @Param chat_id path string true "Chat ID (uuid)"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.ChatResponseSchema
// @Router /chats/{chat_id} [get]
// @Security BearerAuth
func (endpoint Endpoint) RetrieveMessages(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)
	// Parse the UUID parameter
	chatID, err := utils.ParseUUID(c.Params("chat_id"))
	if err != nil {
		return c.Status(400).JSON(err)
	}
	chat := chatManager.GetSingleUserChatFullDetails(db, user, *chatID)
	if chat == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no chat with that ID"))
	}
	
	// Paginate, Convert type and return Messages
	convertedChat := utils.ConvertStructData(chat, schemas.ChatSchema{}).(*schemas.ChatSchema)
	paginatedData, paginatedMessages, err := PaginateQueryset(chat.Edges.Messages, c, 400)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	convertedMessages := utils.ConvertStructData(paginatedMessages, []schemas.MessageSchema{}).(*[]schemas.MessageSchema)
	response := schemas.ChatResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Messages fetched"}.Init(),
		Data: schemas.MessagesSchema{
			Chat: *convertedChat,
			Messages: schemas.MessagesResponseDataSchema{
				PaginatedResponseDataSchema: *paginatedData,
				Items:                       *convertedMessages,
			}.Init(),
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update a Group Chat
// @Description `This endpoint updates a group chat.`
// @Tags Chat
// @Param chat_id path string true "Chat ID (uuid)"
// @Param chat body schemas.GroupChatInputSchema true "Chat object"
// @Success 200 {object} schemas.GroupChatInputResponseSchema
// @Router /chats/{chat_id} [patch]
// @Security BearerAuth
func (endpoint Endpoint) UpdateGroupChat(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)

	chatID, err := utils.ParseUUID(c.Params("chat_id"))
	if err != nil {
		return c.Status(400).JSON(err)
	}

	chatData := schemas.GroupChatInputSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &chatData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(chatData); err != nil {
		return c.Status(422).JSON(err)
	}

	chat := chatManager.GetUserGroup(db, user, *chatID, true)
	if chat == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User owns no group chat with that ID"))
	}
	var errData *utils.ErrorResponse
	chat, errData = chatManager.UpdateGroup(db, chat, chatData)
	if errData != nil {
		return c.Status(422).JSON(errData)
	}
	// Convert type and return chat
	convertedChat := utils.ConvertStructData(chat, schemas.GroupChatInputResponseDataSchema{}).(*schemas.GroupChatInputResponseDataSchema)
	response := schemas.GroupChatInputResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Chat updated"}.Init(),
		Data:           convertedChat.Init(chatData.FileType),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete a Group Chat
// @Description `This endpoint deletes a group chat.`
// @Tags Chat
// @Param chat_id path string true "Chat ID (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Router /chats/{chat_id} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeleteGroupChat(c *fiber.Ctx) error {
	db := endpoint.DB
	chatID, err := utils.ParseUUID(c.Params("chat_id"))
	if err != nil {
		return c.Status(400).JSON(err)
	}
	user := c.Locals("user").(*ent.User)

	// Retrieve & Validate Chat Existence
	chat := chatManager.GetUserGroup(db, user, *chatID)
	if chat == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User owns no group chat with that ID"))
	}
	// Delete and return response
	db.Chat.DeleteOne(chat).Exec(managers.Ctx)
	response := schemas.ResponseSchema{Message: "Group Chat Deleted"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Update a message
// @Description `This endpoint updates a message.`
// @Description
// @Description `You must either send a text or a file or both.`
// @Description
// @Description `The file_upload_data in the response is what is used for uploading the file to cloudinary from client.`
// @Tags Chat
// @Param message_id path string true "Message ID (uuid)"
// @Param message body schemas.MessageUpdateSchema true "Message object"
// @Success 200 {object} schemas.MessageCreateResponseSchema
// @Router /chats/messages/{message_id} [put]
// @Security BearerAuth
func (endpoint Endpoint) UpdateMessage(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)

	messageID, err := utils.ParseUUID(c.Params("message_id"))
	if err != nil {
		return c.Status(400).JSON(err)
	}

	message := messageManager.GetUserMessage(db, user, *messageID)
	if message == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no message with that ID"))
	}

	messageData := schemas.MessageUpdateSchema{}
	// Validate request
	if errCode, errData := DecodeJSONBody(c, &messageData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(messageData); err != nil {
		return c.Status(422).JSON(err)
	}

	message = messageManager.Update(db, message, messageData.Text, messageData.FileType)
	// Convert type and return message
	convertedMessage := utils.ConvertStructData(message, schemas.MessageCreateResponseDataSchema{}).(*schemas.MessageCreateResponseDataSchema)
	response := schemas.MessageCreateResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Message updated"}.Init(),
		Data:           convertedMessage.Init(messageData.FileType),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete a message
// @Description `This endpoint deletes a message.`
// @Tags Chat
// @Param message_id path string true "Message ID (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Router /chats/messages/{message_id} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeleteMessage(c *fiber.Ctx) error {
	db := endpoint.DB
	messageID, err := utils.ParseUUID(c.Params("message_id"))
	if err != nil {
		return c.Status(400).JSON(err)
	}
	user := c.Locals("user").(*ent.User)

	// Retrieve & Validate Message Existence
	message := messageManager.GetUserMessage(db, user, *messageID)
	if message == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no message with that ID"))
	}
	chat := message.Edges.Chat
	messagesCount := chatManager.GetMessagesCount(db, chat.ID)

	// Send message deletion socket
	SendMessageDeletionInSocket(c, chat.ID, message.ID)

	// Delete message and chat if its the last message in the dm being deleted
	if messagesCount == 1 && chat.Ctype == "DM" {
		db.Chat.DeleteOne(chat).Exec(managers.Ctx) // Message deletes if chat gets deleted (CASCADE)
	} else {
		db.Message.DeleteOne(message).Exec(managers.Ctx)
	}

	// Return response
	response := schemas.ResponseSchema{Message: "Message Deleted"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Create a Group Chat
// @Description `This endpoint creates a group chat.`
// @Description
// @Description `The users_entry field should be a list of usernames you want to add to the group.`
// @Description
// @Description `Note: You cannot add more than 99 users in a group (1 owner + 99 other users = 100 users total).`
// @Tags Chat
// @Param chat body schemas.GroupChatCreateSchema true "Chat object"
// @Success 201 {object} schemas.GroupChatInputResponseSchema
// @Router /chats/groups/group [post]
// @Security BearerAuth
func (endpoint Endpoint) CreateGroupChat(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*ent.User)

	chatData := schemas.GroupChatCreateSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &chatData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(chatData); err != nil {
		return c.Status(422).JSON(err)
	}
	usersToAdd := userManager.GetByUsernames(db, chatData.UsernamesToAdd, user.ID)
	if len(usersToAdd) == 0 {
		data := map[string]string{
			"usernames_to_add": "Enter at least one valid username",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}
	chat := chatManager.CreateGroup(db, user, usersToAdd, chatData)
	// Convert type and return chat
	convertedChat := utils.ConvertStructData(chat, schemas.GroupChatInputResponseDataSchema{}).(*schemas.GroupChatInputResponseDataSchema)
	response := schemas.GroupChatInputResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Chat created"}.Init(),
		Data:           convertedChat.Init(chatData.FileType),
	}
	return c.Status(201).JSON(response)
}
