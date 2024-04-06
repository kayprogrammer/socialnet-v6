package schemas

import (
	"time"

	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
)

type ProfileUpdateSchema struct {
	FirstName *string    `json:"first_name" validate:"omitempty,max=50,min=1" example:"John"`
	LastName  *string    `json:"last_name" validate:"omitempty,max=50,min=1" example:"Doe"`
	Bio       *string    `json:"bio" validate:"omitempty,max=200" example:"Software Engineer | Go Fiber Developer"`
	Dob       *time.Time `json:"dob" validate:"omitempty" example:"2001-01-16T00:00:00.106416+01:00"`
	CityID    *uuid.UUID `json:"city_id" validate:"omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	FileType  *string    `json:"file_type" example:"image/jpeg" validate:"omitempty,file_type_validator"`
}

func (p ProfileUpdateSchema) SetValues(user *models.User) *models.User {
	if p.FirstName != nil {
		user.FirstName = *p.FirstName
	}
	if p.LastName != nil {
		user.LastName = *p.LastName
	}
	user.Bio = p.Bio
	user.Dob = p.Dob
	return user
}

type DeleteUserSchema struct {
	Password string `json:"password" validate:"required" example:"password"`
}

type SendFriendRequestSchema struct {
	Username string `json:"username" validate:"required" example:"john-doe"`
}

type AcceptFriendRequestSchema struct {
	SendFriendRequestSchema
	Accepted bool `json:"accepted" example:"true"`
}

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



type ReadNotificationSchema struct {
	MarkAllAsRead bool       `json:"mark_all_as_read" example:"false"`
	ID            *uuid.UUID `json:"id" validate:"required_if=MarkAllAsRead false,omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
}

// RESPONSE SCHEMAS
// CITIES
type CitiesResponseSchema struct {
	ResponseSchema
	Data []models.City `json:"data"`
}

func (data CitiesResponseSchema) Init() CitiesResponseSchema {
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
	Items []models.User `json:"users"`
}

func (data ProfilesResponseDataSchema) Init() ProfilesResponseDataSchema {
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
	Data ProfilesResponseDataSchema `json:"data"`
}

type ProfileResponseSchema struct {
	ResponseSchema
	Data models.User `json:"data"`
}

type ProfileUpdateResponseDataSchema struct {
	models.User
	FileUploadData *utils.SignatureFormat `json:"file_upload_data"`
}

func (profileData ProfileUpdateResponseDataSchema) Init(fileType *string) ProfileUpdateResponseDataSchema {
	image := profileData.User.AvatarObj
	if fileType != nil && image != nil { // Generate data when file is being uploaded
		fuData := utils.GenerateFileSignature(image.ID.String(), "avatars")
		profileData.FileUploadData = &fuData
	}
	profileData.User = profileData.User.Init()
	return profileData
}

type ProfileUpdateResponseSchema struct {
	ResponseSchema
	Data ProfileUpdateResponseDataSchema `json:"data"`
}

// NOTIFICATIONS
type NotificationsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []models.Notification `json:"notifications"`
}

func (data NotificationsResponseDataSchema) Init(currentUserID uuid.UUID) NotificationsResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init(currentUserID)
	}
	data.Items = items
	return data
}

type NotificationsResponseSchema struct {
	ResponseSchema
	Data NotificationsResponseDataSchema `json:"data"`
}
