package routes

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

// Entry & Exit Schemas
type SocketMessageEntrySchema struct {
	Status string    `json:"status" validate:"required,oneof=CREATED UPDATED DELETED"`
	ID     uuid.UUID `json:"id" validate:"required"`
}

type SocketMessageExitSchema struct {
	models.Message
	Status string `json:"status"`
}

// ---------------------------

var (
	messageData    = SocketMessageEntrySchema{}
)

// Retrieve chat or user based on the given id
func GetChatOrUser(c *websocket.Conn, db *gorm.DB, user models.User, id string) (models.Chat, models.User) {
	chat := models.Chat{}
	objUser := models.User{}

	if user.ID.String() != id {
		parsedID, _ := utils.ParseUUID(id)
		if parsedID == nil {
			db.Where(models.User{Username: id}).Take(&objUser, objUser)
		} else {
			chat = chatManager.GetByID(db, *parsedID)
		}
	} else {
		objUser = user // Message is sent to self
	}
	c.Locals("objUser", objUser)
	return chat, objUser
}

// --------------------------------------------

// Validate chat existence or membership
func ValidateChatMembership(c *websocket.Conn, db *gorm.DB, user models.User, id string) (*int, *string, *string) {
	chat, objUser := GetChatOrUser(c, db, user, id)
	if chat.ID == nil && objUser.ID == nil {
		// If no chat nor user
		errCode := 4004
		errType := "invalid_input"
		errMsg := "Invalid ID"
		return &errCode, &errType, &errMsg
	}
	if chat.ID != nil && user.ID.String() != chat.OwnerID.String() && !chatManager.UserIsMember(chat, user) {
		errCode := 4001
		errType := "invalid_member"
		errMsg := "You're not a member of this chat"
		return &errCode, &errType, &errMsg
	}
	return nil, nil, nil
}

// ----------------------------------------------------

// Store new connection client
func AddChatClient(c *websocket.Conn, db *gorm.DB, id string) (*int, *string, *string) {
	// Validate chat ID & membership
	user := c.Locals("user").(*models.User)
	secret := c.Locals("secret").(*string)

	if secret == nil {
		// validate chat memership
		errCode, errType, errMsg := ValidateChatMembership(c, db, *user, id)
		if errCode != nil && errType != nil && errMsg != nil {
			return errCode, errType, errMsg
		}
	}
	// Add client
	AddClient(c)
	return nil, nil, nil
}

// ------------------------------------------

// Validate data entering the socket.
func ValidateEnteredData(c *websocket.Conn, db *gorm.DB, user *models.User, secret *string, data []byte) (*[]byte, *int, *string, *string, *map[string]string) {
	// Ensure data is a Message data. That means it aligns with the Message schema above
	err := json.Unmarshal(data, &messageData)
	if err != nil {
		errCode := 4220
		errType := utils.ERR_INVALID_ENTRY
		errMsg := "Invalid Json data"
		return nil, &errCode, &errType, &errMsg, nil
	} else if err := validator.Validate(messageData); err != nil {
		errCode := 4220
		errType := err.Code
		errMsg := "Invalid Message data"
		errData := err.Data
		return nil, &errCode, &errType, &errMsg, errData
	}
	status := messageData.Status
	if status == "DELETED" && secret == nil {
		// Only allowed for secret users (in app)
		errCode := 4001
		errType := utils.ERR_UNAUTHORIZED_USER
		errMsg := "Not allowed to send deletion socket message"
		return nil, &errCode, &errType, &errMsg, nil
	}
	messageDataToReturn := data
	if status != "DELETED" {
		message := messageManager.GetByID(db, messageData.ID)
		if message.ID == nil {
			errCode := 4004
			errType := utils.ERR_NON_EXISTENT
			errMsg := "Invalid message ID"
			return nil, &errCode, &errType, &errMsg, nil
		} else if message.SenderID.String() != user.ID.String() {
			errCode := 4001
			errType := utils.ERR_INVALID_OWNER
			errMsg := "Message isn't yours"
			return nil, &errCode, &errType, &errMsg, nil
		}
		messageData := SocketMessageExitSchema{
			Message: message.Init(),
			Status:  messageData.Status,
		}
		messageDataJson, _ := json.Marshal(messageData)
		messageDataToReturn = []byte(messageDataJson)
	}
	return &messageDataToReturn, nil, nil, nil, nil
}

// --------------------------------------------

// Broadcast chat messages to connected clients
func broadcastChatMessage(c *websocket.Conn, mt int, groupName string, data []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		secret := client.Locals("secret").(*string)

		// Only true receivers should access the data
		if client.Locals("groupName") == groupName && secret == nil {
			user := client.Locals("user").(*models.User)
			objUser := client.Locals("objUser").(*models.User)
			if objUser != nil {
				// Ensure that reading messages from a user id can only be done by the owner
				if user.ID.String() == objUser.ID.String() {
					if err := client.WriteMessage(mt, data); err != nil {
						log.Println("write:", err)
					}
				}
			} else {
				if err := client.WriteMessage(mt, data); err != nil {
					log.Println("write:", err)
				}
			}
		}
	}
}

// --------------------------------------------

// Chat socket endpoint
func ChatSocket(c *websocket.Conn) {
	db := database.ConnectDb(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	token := c.Headers("Authorization")
	chatID := c.Params("id")

	var (
		mt        int
		entryData []byte
		exitData  *[]byte
		err       error
		user      *models.User
		secret    *string
		errC      *int               // error code
		errT      *string            // error type
		errM      *string            // error message
		errD      *map[string]string // error data
	)

	// Validate Auth
	if user, secret, errM = ValidateAuth(db, token); errM != nil {
		ReturnError(c, utils.ERR_INVALID_TOKEN, *errM, 4001)
		return
	}
	c.Locals("user", user)
	c.Locals("secret", secret)

	// Set Group name
	groupName := "chat_" + chatID
	c.Locals("groupName", groupName)

	// Add the client to the list of connected clients
	errC, errT, errM = AddChatClient(c, db, chatID)
	if errC != nil || errT != nil || errM != nil {
		ReturnError(c, *errT, *errM, *errC)
		return
	}

	// Remove the client from the list when the handler exits
	defer RemoveClient(c)

	for {
		if mt, entryData, err = c.ReadMessage(); err != nil {
			ReturnError(c, utils.ERR_INVALID_ENTRY, "Invalid Entry", 4220)
			break
		}

		// Validate received data
		exitData, errC, errT, errM, errD = ValidateEnteredData(c, db, user, secret, entryData)
		if errC != nil {
			ReturnError(c, *errT, *errM, *errC, errD)
			break
		}
		broadcastChatMessage(c, mt, groupName, *exitData)
	}
}

// --------------------------------------------
