package managers

import (
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ----------------------------------
// FRIEND MANAGEMENT
// --------------------------------
type FriendManager struct {
}

func (obj FriendManager) GetFriends(db *gorm.DB, user models.User) []models.User {
	friends := []models.Friend{}
	db.Where(
		models.Friend{Status: choices.FACCEPTED},
	).Where(
		models.Friend{RequesterID: user.ID}).Or(models.Friend{RequesteeID: user.ID}).Find(&friends)

	friendIDs := []uuid.UUID{}
	for i := range friends {
		requesterID := friends[i].RequesterID
		requesteeID := friends[i].RequesteeID
		if user.ID.String() == requesterID.String() {
			friendIDs = append(friendIDs, requesteeID)
		} else {
			friendIDs = append(friendIDs, requesterID)
		}
	}
	users := []models.User{}
	if len(friendIDs) > 0 {
		db.Preload(clause.Associations).Find(&users, friendIDs)
	}
	return users
}

func (obj FriendManager) GetFriendRequests(db *gorm.DB, user *models.User) []models.User {
	friendObjects := []models.Friend{}
	db.Select("requester_id").Where(models.Friend{RequesteeID: user.ID, Status: choices.FPENDING}).Find(&friendObjects)

	friendIDs := []uuid.UUID{}
	for i := range friendObjects {
		friendIDs = append(friendIDs, friendObjects[i].RequesterID)
	}

	friends := []models.User{}
	if len(friendIDs) > 0 {
		db.Preload(clause.Associations).Find(&friends, friendIDs)
	}
	return friends
}

func (obj FriendManager) GetRequesteeAndFriendObj(db *gorm.DB, user *models.User, username string, statusOpts ...choices.FriendStatusChoice) (*models.User, *models.Friend, *utils.ErrorResponse) {
	requestee := models.User{Username: username}
	db.Take(&requestee, requestee)
	if requestee.ID == nil {
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "User does not exist!")
		return nil, nil, &errData
	}
	friend := models.Friend{}
	fq := db.Where(models.Friend{RequesterID: user.ID, RequesteeID: requestee.ID}).Or(models.Friend{RequesterID: requestee.ID, RequesteeID: user.ID})
	if len(statusOpts) > 0 {
		// If status param is provided
		fq = fq.Where(models.Friend{Status: statusOpts[0]})
	}
	fq.Take(&friend, friend)
	return &requestee, &friend, nil
}

// func (obj FriendManager) DropData(db *gorm.DB) {
// 	db.Friend.Delete().ExecX(Ctx)
// }

// ----------------------------------
// NOTIFICATION MANAGEMENT
// --------------------------------
type NotificationManager struct {
}

func (obj NotificationManager) GetQueryset(db *gorm.DB, userID uuid.UUID) []models.Notification {
	notifications := []models.Notification{}
	db.Preload(clause.Associations).Order("created_at DESC").Find(&notifications)
	return notifications
}

func (obj NotificationManager) MarkAsRead(db *gorm.DB, user *models.User) {
	notifications := []models.Notification{}
	db.Model(&user).Find(&notifications, "NotificationsReceived")
	db.Model(&user).Association("NotificationsRead").Append(&notifications)
}

func (obj NotificationManager) ReadOne(db *gorm.DB, user *models.User, notificationID uuid.UUID) *utils.ErrorResponse {
	notification := models.Notification{}
	db.Model(&user).Take(&notification, notificationID, "NotificationsReceived")

	if notification.ID == nil {
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no notification with that ID")
		return &errData
	}
	db.Model(&user).Association("NotificationsRead").Append(&notification)
	return nil
}

func (obj NotificationManager) Create(db *gorm.DB, sender *models.User, ntype choices.NotificationChoice, receivers []models.User, post *models.Post, comment *models.Comment, reply *models.Reply, text *string) models.Notification {
	// Create Notification
	notification := models.Notification{Ntype: ntype, Text: text, Sender: sender, Post: post, Comment: comment, Reply: reply, Receivers: receivers}
	db.Create(&notification)
	return notification
}

func (obj NotificationManager) GetOrCreate(db *gorm.DB, sender *models.User, ntype choices.NotificationChoice, receivers []models.User, post *models.Post, comment *models.Comment, reply *models.Reply) (models.Notification, bool) {
	created := false
	notification := models.Notification{Sender: sender, Ntype: ntype, Post: post, Comment: comment, Reply: reply}
	db.Preload("Sender", "Post", "Comment", "Reply").Take(&notification, notification)
	if notification.ID == nil {
		created = true
		// Create notification
		notification = obj.Create(db, sender, ntype, receivers, post, comment, reply, nil)
	}
	return notification, created
}

func (obj NotificationManager) Get(db *gorm.DB, sender *models.User, ntype choices.NotificationChoice, post *models.Post, comment *models.Comment, reply *models.Reply) *models.Notification {
	notification := models.Notification{Sender: sender, Ntype: ntype, Post: post, Comment: comment, Reply: reply}
	db.Take(&notification, notification)
	if notification.ID == nil {
		return nil
	}
	return &notification
}

func (obj NotificationManager) IsAmongReceivers(db *gorm.DB, notificationID uuid.UUID, receiverID uuid.UUID) bool {
	notification := models.Notification{}
	db.Preload("Receivers").Take(&notification, notificationID)
	if notification.ID == nil {
		return false
	}

	// Check if user in notification receivers
	found := false
	for _, item := range notification.Receivers {
		if item.ID.String() == receiverID.String() {
			found = true
			break
		}
	}
	return found
}

// func (obj NotificationManager) DropData(db *gorm.DB) {
// 	db.Notification.Delete().ExecX(Ctx)
// }
