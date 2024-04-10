package routes

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/managers"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"gorm.io/gorm"
)

type SocketNotificationSchema struct {
	models.Notification
	Status string `json:"status"`
}

var notificationObj SocketNotificationSchema

// Function to broadcast a notification data to all connected clients
func broadcastNotificationMessage(db *gorm.DB, mt int, msg []byte) {
	notificationManager := managers.NotificationManager{}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		user := client.Locals("user").(*models.User)
		if user == nil {
			continue
		}
		json.Unmarshal(msg, &notificationObj)
		// Ensure user is a valid recipient of this notification
		userIsAmongReceiver := notificationManager.IsAmongReceivers(db, notificationObj.ID, user.ID)
		if userIsAmongReceiver {
			if err := client.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
			}
		}
	}
	// Delete comment or reply here after the socket message has been sent for comment & reply deletion
	// Although another better way will be to delete the comment or reply the respective view/handler
	// But then the notification will be deleted alongside (cos of CASCADE relationship) before the notification socket will be sent
	// Which will prevent the user from seeing the real time notification cos the IsAmongReceivers won't work with an already deleted notifiation
	// To prevent this you can just set the relationship to SetNull, then delete notification here, and delete comment & reply in the view.
	// The only drawback I can think of concerning the below method is that if by any means there was an issue with the socket, the stuff won't get deleted (will probably implement a better solution in another version of this project).
	// Omo na wahala be that oh. But anyway, just go ahead with the SetNull whatever. I'm too lazy to change anything now.
	// Sorry for the long note (no vex)
	if notificationObj.Status == "DELETED" && notificationObj.Ntype != choices.NREACTION {
		db.Delete(&notificationObj.Notification)
		if notificationObj.CommentSlug != nil {
			var commentSlug string = *notificationObj.CommentSlug
			db.Delete(&models.Comment{}, "slug = ?", commentSlug)
		} else if notificationObj.ReplySlug != nil {
			var replySlug string = *notificationObj.ReplySlug
			db.Delete(&models.Reply{}, "slug = ?", replySlug)
		}
	}
}

func (ep Endpoint) NotificationSocket(c *websocket.Conn) {
	db := database.ConnectDb(cfg, true)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	token := c.Headers("Authorization")

	var (
		mt     int
		msg    []byte
		err    error
		user   *models.User
		secret *string
		errM   *string
	)

	// Validate Auth
	if user, secret, errM = ValidateAuth(db, token); errM != nil {
		ReturnError(c, utils.ERR_INVALID_TOKEN, *errM, 4001)
		return
	}
	// Add the client to the list of connected clients
	c.Locals("user", user)
	AddClient(c)

	// Remove the client from the list when the handler exits
	defer RemoveClient(c)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			ReturnError(c, utils.ERR_INVALID_ENTRY, "Invalid Entry", 4220)
			break
		}

		// Notifications can only be broadcasted from the app using the socket secret
		if secret != nil {
			broadcastNotificationMessage(db, mt, msg)
		} else {
			ReturnError(c, utils.ERR_UNAUTHORIZED_USER, "Not authorized to send data", 4001)
			break
		}
	}
}
