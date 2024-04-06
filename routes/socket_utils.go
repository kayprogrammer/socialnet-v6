package routes

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"gorm.io/gorm"
)

// Maintain db & a list of connected clients
var (
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
	validator    = utils.Validator()
)

// Function to add a client to the list
func AddClient(c *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[c] = true
}

// Function to remove a client from the list
func RemoveClient(c *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, c)
}

type ErrorResp struct {
	Status  string             `json:"status"`
	Code    int                `json:"code"`
	Type    string             `json:"type"`
	Message string             `json:"message"`
	Data    *map[string]string `json:"data,omitempty"`
}

func ReturnError(c *websocket.Conn, errType string, message string, code int, dataOpts ...*map[string]string) {
	errorResponse := ErrorResp{Status: "failure", Code: code, Type: errType, Message: message}
	if len(dataOpts) > 0 {
		errorResponse.Data = dataOpts[0]
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	c.WriteMessage(websocket.TextMessage, jsonResponse)
}

func ValidateAuth(db *gorm.DB, token string) (*models.User, *string, *string) {
	var (
		errMsg *string
		secret *string
		user   *models.User
	)
	if len(token) < 1 {
		err := "Auth bearer not set"
		errMsg = &err
	} else if token == config.GetConfig().SocketSecret {
		secret = &token
	} else {
		// Get User
		userObj, err := GetUser(token, db)
		if err != nil {
			errMsg = err
		}
		user = userObj
	}
	return user, secret, errMsg
}
