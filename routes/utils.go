package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
)

func SuccessResponse(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func ValidateReactionFocus(focus choices.FocusTypeChoice) *utils.ErrorResponse {
	switch focus {
	case "POST", "COMMENT", "REPLY":
		return nil
	}
	err := utils.RequestErr(utils.ERR_INVALID_VALUE, "Invalid 'focus' value")
	return &err
}

func SendNotificationInSocket(fiberCtx *fiber.Ctx, notification models.Notification, commentSlug *string, replySlug *string, statusOpts ...string) error {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return nil
	}
	
	// Check if page size is provided as an argument
	status := "CREATED"
	if len(statusOpts) > 0 {
		status = statusOpts[0]
	}
	webSocketScheme := "ws://"
	if fiberCtx.Secure() {
		webSocketScheme = "wss://"
	}
	uri := webSocketScheme + fiberCtx.Hostname() + "/api/v6/ws/notifications/"
	notificationData := SocketNotificationSchema{
		Notification: models.Notification{BaseModel: models.BaseModel{ID: notification.ID}, Ntype: notification.Ntype, CommentSlug: commentSlug, ReplySlug: replySlug},
		Status:             status,
	}
	if status == "CREATED" {
		notificationData = SocketNotificationSchema{
			Notification: notification.Init(nil),
			Status:             status,
		}
	}

	// Connect to the WebSocket server
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	headers := make(http.Header)
	headers.Add("Authorization", cfg.SocketSecret)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Marshal the notification data to JSON
	data, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}

	// Send the notification to the WebSocket server
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	// Close the WebSocket connection
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func SendMessageDeletionInSocket(fiberCtx *fiber.Ctx, chatID uuid.UUID, messageID uuid.UUID) error {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return nil
	}
	webSocketScheme := "ws://"
	if fiberCtx.Secure() {
		webSocketScheme = "wss://"
	}
	uri := webSocketScheme + fiberCtx.Hostname() + "/api/v6/ws/chats/" + chatID.String()
	chatData := SocketMessageEntrySchema{
		ID:     messageID,
		Status: "DELETED",
	}

	// Connect to the WebSocket server
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	headers := make(http.Header)
	headers.Add("Authorization", cfg.SocketSecret)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Marshal the notification data to JSON
	data, err := json.Marshal(chatData)
	if err != nil {
		return err
	}

	// Send the message to the WebSocket server
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	// Close the WebSocket connection
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
