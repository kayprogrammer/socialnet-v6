package schemas

import (
	"github.com/kayprogrammer/socialnet-v6/models"
)


// type ProfileSchema struct {
// 	Edges        	*ent.UserEdges 		`json:"edges,omitempty" swaggerignore:"true"`
//     FirstName 		string				`json:"first_name" example:"John"`
//     LastName 		string				`json:"last_name" example:"Doe"`
//     Username 		string				`json:"username" example:"john-doe"`
//     Email 			string				`json:"email" example:"johndoe@email.com"`
//     Avatar 			*string				`json:"avatar" example:"https://img.com"`
//     Bio 			*string				`json:"bio" example:"Software Engineer | Go Fiber Developer"`
//     Dob 			*time.Time			`json:"dob"`
//     City 			*string				`json:"city" example:"Lekki"`
//     CreatedAt 		*time.Time			`json:"created_at"`
//     UpdatedAt 		*time.Time			`json:"updated_at"`
// }

// func (user ProfileSchema) Init () ProfileSchema {
// 	// Set Related Data.
// 	avatar := user.Edges.Avatar 
// 	if avatar != nil {
// 		url := utils.GenerateFileUrl(avatar.ID.String(), "avatars", avatar.ResourceType)
// 		user.Avatar = &url
// 	}
// 	city := user.Edges.City
// 	if city != nil {
// 		user.City = &city.Name
// 	}
// 	user.Edges = nil // Omit edges
// 	return user
// }

// type ProfileUpdateSchema struct {
// 	FirstName      *string 		`json:"first_name" validate:"omitempty,max=50,min=1" example:"John"`
// 	LastName       *string 		`json:"last_name" validate:"omitempty,max=50,min=1" example:"Doe"`
// 	Bio			   *string 		`json:"bio" validate:"omitempty,max=200" example:"Software Engineer | Go Fiber Developer"`
// 	Dob			   *time.Time 	`json:"dob" validate:"omitempty" example:"2001-01-16T00:00:00.106416+01:00"`
// 	CityID		   *uuid.UUID	`json:"city_id" validate:"omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// 	FileType	   *string		`json:"file_type" example:"image/jpeg" validate:"omitempty,file_type_validator"`
// 	City		   *ent.City	`json:"-"`
// }

// type DeleteUserSchema struct {
// 	Password 		string		`json:"password" validate:"required" example:"password"`
// }

// type SendFriendRequestSchema struct {
// 	Username		string		`json:"username" validate:"required" example:"john-doe"`
// }

// type AcceptFriendRequestSchema struct {
// 	Username		string		`json:"username" validate:"required" example:"john-doe"`
// 	Accepted		bool		`json:"accepted" example:"true"`
// }

// type NotificationSchema struct {
// 	Edges        	*ent.NotificationEdges 	`json:"edges,omitempty" swaggerignore:"true"`
// 	ID				uuid.UUID				`json:"id"`
// 	Sender			*UserDataSchema			`json:"sender"`
// 	Ntype			string					`json:"ntype" example:"REACTION"`
// 	Text        	*string 				`json:"text,omitempty" swaggerignore:"true"`
// 	Message			string					`json:"message" example:"John Doe reacted to your post"`
// 	PostSlug		*string					`json:"post_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// 	CommentSlug		*string					`json:"comment_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// 	ReplySlug		*string					`json:"reply_slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// 	IsRead			bool					`json:"is_read" example:"true"`
// }

// func (notification NotificationSchema) Init (currentUserID *uuid.UUID) NotificationSchema {
// 	// Set Related Data.
// 	sender := notification.Edges.Sender
// 	if sender != nil {
// 		senderData := UserDataSchema{}.Init(sender)
// 		notification.Sender = &senderData
// 	}

// 	// Set Target slug
// 	notification = notification.SetTargetSlug()
// 	// Set Notification message
// 	text := notification.Text
// 	if text == nil {
// 		notificationMsg := notification.GetMessage()
// 		text = &notificationMsg
// 	}
// 	notification.Message = *text 
// 	notification.Text = nil // Omit text

// 	// Set IsRead
// 	if currentUserID != nil {
// 		readBy := notification.Edges.ReadBy
// 		for _, user := range readBy {
// 			if user.ID == *currentUserID {
// 				notification.IsRead = true
// 				break
// 			}
// 		}
// 	}
// 	notification.Edges = nil // Omit edges
// 	return notification
// }

// func (notification NotificationSchema) SetTargetSlug() NotificationSchema {
// 	post := notification.Edges.Post
// 	comment := notification.Edges.Comment
// 	reply := notification.Edges.Reply
// 	if post != nil {
// 		notification.PostSlug = &post.Slug
// 	} else if comment != nil {
// 		notification.CommentSlug = &comment.Slug
// 	} else if reply != nil {
// 		notification.ReplySlug = &reply.Slug
// 	}
// 	return notification

// }

// func (notification NotificationSchema) GetMessage() string {
// 	ntype := notification.Ntype
// 	sender := notification.Sender.Name
// 	message := sender + " reacted to your post"
// 	if ntype == "REACTION" {
// 		if notification.CommentSlug != nil {
// 			message = sender + " reacted to your comment"
// 		} else if notification.ReplySlug != nil {
// 			message = sender + " reacted to your reply"
// 		}
// 	} else if ntype == "COMMENT" {
// 		message = sender + " commented on your post"
// 	} else if ntype == "REPLY" {
// 		message = sender + " replied your comment"
// 	}
// 	return message
// }

// type ReadNotificationSchema struct {
// 	MarkAllAsRead		bool			`json:"mark_all_as_read" example:"false"`
// 	ID					*uuid.UUID		`json:"id" validate:"required_if=MarkAllAsRead false,omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// }

// RESPONSE SCHEMAS
// CITIES
type CitiesResponseSchema struct {
	ResponseSchema
	Data			[]models.City		`json:"data"`
}

func (data CitiesResponseSchema) Init () CitiesResponseSchema {
	// Set Initial Data
	cities := data.Data
	for i := range cities {
		cities[i] = cities[i].Init()
	}
	data.Data = cities
	return data
}

// USERS
type ProfilesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]models.User		`json:"users"`
}

func (data ProfilesResponseDataSchema) Init () ProfilesResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type ProfilesResponseSchema struct {
	ResponseSchema
	Data			ProfilesResponseDataSchema		`json:"data"`
}

type ProfileResponseSchema struct {
	ResponseSchema
	Data			models.User		`json:"data"`
}

// type ProfileUpdateResponseDataSchema struct {
// 	ProfileSchema
// 	FileUploadData 		*utils.SignatureFormat 	`json:"file_upload_data"`
// }

// func (profileData ProfileUpdateResponseDataSchema) Init(fileType *string) ProfileUpdateResponseDataSchema {
// 	image := profileData.ProfileSchema.Edges.Avatar
// 	if fileType != nil && image != nil { // Generate data when file is being uploaded
// 		fuData := utils.GenerateFileSignature(image.ID.String(), "avatars")
// 		profileData.FileUploadData = &fuData
// 	}
// 	profileData.ProfileSchema = profileData.ProfileSchema.Init()
// 	return profileData	
// }

// type ProfileUpdateResponseSchema struct {
// 	ResponseSchema
// 	Data				ProfileUpdateResponseDataSchema			`json:"data"`
// }

// // NOTIFICATIONS
// type NotificationsResponseDataSchema struct {
// 	PaginatedResponseDataSchema
// 	Items			[]NotificationSchema		`json:"notifications"`
// }

// func (data NotificationsResponseDataSchema) Init (currentUserID uuid.UUID) NotificationsResponseDataSchema {
// 	// Set Initial Data
// 	items := data.Items
// 	for i := range items {
// 		items[i] = items[i].Init(&currentUserID)
// 	}
// 	data.Items = items
// 	return data
// }

// type NotificationsResponseSchema struct {
// 	ResponseSchema
// 	Data			NotificationsResponseDataSchema		`json:"data"`
// }