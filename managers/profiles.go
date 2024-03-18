package managers

import (
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
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
		models.Friend{RequesterID: user.ID}).Or(models.Friend{RequesteeID: user.ID},
	).Find(&friends)

	friendIDs := []uuid.UUID{}
	for i := range friends {
		requesterID := friends[i].RequesterID
		requesteeID := friends[i].RequesteeID
		if user.ID.String() == requesterID.String()  {
			friendIDs = append(friendIDs, requesteeID)
		} else {
			friendIDs = append(friendIDs, requesterID)
		}
	}
	users := []models.User{}
	db.Preload(clause.Associations).Find(&users, friendIDs)
	return users
}

// func (obj FriendManager) GetFriendRequests(client *ent.Client, userObj *ent.User) []*ent.User {
// 	friendObjects, _ := client.Friend.Query().
// 		Where(
// 			friend.RequesteeIDEQ(userObj.ID),
// 			friend.StatusEQ("PENDING"),
// 		).
// 		All(Ctx)
// 	var friendIDs []uuid.UUID
// 	for i := range friendObjects {
// 		friendIDs = append(friendIDs, friendObjects[i].RequesterID)
// 	}

// 	friends, _ := client.User.Query().
// 		Where(user.IDIn(friendIDs...)).
// 		WithCity().
// 		WithAvatar().
// 		All(Ctx)
// 	return friends
// }

// func (obj FriendManager) GetRequesteeAndFriendObj(client *ent.Client, userObj *ent.User, username string, statusOpts ...string) (*ent.User, *ent.Friend, *utils.ErrorResponse) {
// 	requestee, _ := client.User.Query().
// 		Where(user.Username(username)).
// 		Only(Ctx)

// 	if requestee == nil {
// 		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "User does not exist!")
// 		return nil, nil, &errData
// 	}
// 	fq := client.Friend.Query().
// 		Where(
// 			friend.Or(
// 				friend.And(
// 					friend.RequesterIDEQ(userObj.ID),
// 					friend.RequesteeIDEQ(requestee.ID),
// 				),
// 				friend.And(
// 					friend.RequesterIDEQ(requestee.ID),
// 					friend.RequesteeIDEQ(userObj.ID),
// 				),
// 			),
// 		)
// 	if len(statusOpts) > 0 {
// 		// If status param is provided
// 		fq = fq.Where(friend.StatusEQ(friend.Status(statusOpts[0])))
// 	}

// 	friend, _ := fq.Only(Ctx)
// 	return requestee, friend, nil
// }

// func (obj FriendManager) Create(client *ent.Client, requester *ent.User, requestee *ent.User, status friend.Status) *ent.Friend {
// 	friendObj := client.Friend.
// 		Create().
// 		SetStatus(status).
// 		SetRequester(requester).
// 		SetRequestee(requestee).
// 		SaveX(Ctx)

// 	// Set related data
// 	friendObj.Edges.Requester = requester
// 	friendObj.Edges.Requestee = requestee
// 	return friendObj
// }

// func (obj FriendManager) DropData(client *ent.Client) {
// 	client.Friend.Delete().ExecX(Ctx)
// }

// // ----------------------------------
// // NOTIFICATION MANAGEMENT
// // --------------------------------
// type NotificationManager struct {
// }

// func (obj NotificationManager) GetQueryset(client *ent.Client, userID uuid.UUID) []*ent.Notification {
// 	notifications := client.Notification.Query().
// 		Where(notification.HasReceiversWith(user.ID(userID))).
// 		WithSender(func(uq *ent.UserQuery) { uq.WithAvatar() }).
// 		WithPost().
// 		WithComment().
// 		WithReply().
// 		WithReadBy().
// 		Order(notification.ByCreatedAt(sql.OrderDesc())).
// 		AllX(Ctx)
// 	return notifications
// }

// func (obj NotificationManager) MarkAsRead(client *ent.Client, userID uuid.UUID) {
// 	client.Notification.
// 		Update().
// 		Where(notification.HasReceiversWith(user.ID(userID))).
// 		AddReadByIDs(userID).
// 		SaveX(Ctx)
// }

// func (obj NotificationManager) ReadOne(client *ent.Client, userID uuid.UUID, notificationID uuid.UUID) *utils.ErrorResponse {
// 	n, _ := client.Notification.
// 		Query().
// 		Where(notification.HasReceiversWith(user.ID(userID)), notification.ID(notificationID)).
// 		Only(Ctx)
// 	if n == nil {
// 		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no notification with that ID")
// 		return &errData
// 	}
// 	n.Update().AddReadByIDs(userID).SaveX(Ctx)
// 	return nil
// }

// func (obj NotificationManager) Create(client *ent.Client, sender *ent.User, ntype notification.Ntype, receiverIDs []uuid.UUID, post *ent.Post, comment *ent.Comment, reply *ent.Reply, text *string) *ent.Notification {
// 	// Create Notification
// 	nc := client.Notification.Create().
// 		SetNtype(ntype).
// 		SetNillableText(text).
// 		AddReceiverIDs(receiverIDs...)

// 	if sender != nil {
// 		nc = nc.SetSender(sender)
// 	}
// 	if post != nil {
// 		nc = nc.SetPost(post)
// 	} else if comment != nil {
// 		nc = nc.SetComment(comment)
// 	} else if reply != nil {
// 		nc = nc.SetReply(reply)
// 	}

// 	notification := nc.SaveX(Ctx)

// 	// Set related data
// 	notification.Edges.Sender = sender
// 	notification.Edges.Post = post
// 	notification.Edges.Comment = comment
// 	notification.Edges.Reply = reply
// 	return notification
// }

// func (obj NotificationManager) GetOrCreate(client *ent.Client, sender *ent.User, ntype notification.Ntype, receiverIDs []uuid.UUID, post *ent.Post, comment *ent.Comment, reply *ent.Reply) (*ent.Notification, bool) {
// 	created := false
// 	nq := client.Notification.Query().
// 		Where(
// 			notification.SenderID(sender.ID),
// 			notification.NtypeEQ(ntype),
// 		)
// 	if post != nil {
// 		nq = nq.Where(notification.PostID(post.ID))
// 	} else if comment != nil {
// 		nq = nq.Where(notification.CommentID(comment.ID))
// 	} else if reply != nil {
// 		nq = nq.Where(notification.ReplyID(reply.ID))
// 	}
// 	n, _ := nq.WithSender().WithPost().WithComment().WithReply().Only(Ctx)
// 	if n == nil {
// 		created = true
// 		// Create notification
// 		n = obj.Create(client, sender, ntype, receiverIDs, post, comment, reply, nil)
// 	}
// 	return n, created
// }

// func (obj NotificationManager) Get(client *ent.Client, sender *ent.User, ntype notification.Ntype, post *ent.Post, comment *ent.Comment, reply *ent.Reply) *ent.Notification {
// 	nq := client.Notification.Query().
// 		Where(notification.SenderID(sender.ID), notification.NtypeEQ(ntype))

// 	if post != nil {
// 		nq = nq.Where(notification.PostID(post.ID))
// 	} else if comment != nil {
// 		nq = nq.Where(notification.CommentID(comment.ID))
// 	} else if reply != nil {
// 		nq = nq.Where(notification.ReplyID(reply.ID))
// 	}
// 	n, _ := nq.Only(Ctx)
// 	return n
// }

// func (obj NotificationManager) IsAmongReceivers(client *ent.Client, notificationID uuid.UUID, receiverID uuid.UUID) bool {
// 	exists := client.Notification.Query().
// 		Where(
// 			notification.ID(notificationID),
// 			notification.HasReceiversWith(user.ID(receiverID)),
// 		).
// 		ExistX(Ctx)
// 	return exists
// }

// func (obj NotificationManager) DropData(client *ent.Client) {
// 	client.Notification.Delete().ExecX(Ctx)
// }
