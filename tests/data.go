package tests

import (
	"fmt"
	"strings"
	"time"

	"github.com/kayprogrammer/socialnet-v6/managers"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/routes"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"gorm.io/gorm"
)

var (
	friendManager       = managers.FriendManager{}
	notificationManager = managers.NotificationManager{}
	chatManager         = managers.ChatManager{}
	messageManager      = managers.MessageManager{}
	postManager         = managers.PostManager{}
	reactionManager     = managers.ReactionManager{}
	commentManager      = managers.CommentManager{}
	replyManager        = managers.ReplyManager{}
)

// AUTH FIXTURES
func CreateTestUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName:      "Test",
		LastName:       "User",
		Email:          "testuser@example.com",
		Password:       "testpassword",
		TermsAgreement: false,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateTestVerifiedUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "Test",
		LastName:  "Verified",
		Email:     "testverifieduser@example.com",
		Password:  "testpassword",
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateAnotherTestVerifiedUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "AnotherTest",
		LastName:  "UserVerified",
		Email:     "anothertestverifieduser@example.com",
		Password:  "testpassword",
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateJwt(db *gorm.DB, user models.User) models.User {
	access := routes.GenerateAccessToken(user.ID, user.Username)
	refresh := routes.GenerateRefreshToken()
	user.Access = &access
	user.Refresh = &refresh
	db.Save(&user)
	return user
}

func AccessToken(db *gorm.DB) string {
	user := CreateTestVerifiedUser(db)
	user = CreateJwt(db, user)
	return *user.Access
}

func AnotherAccessToken(db *gorm.DB) string {
	user := CreateAnotherTestVerifiedUser(db)
	user = CreateJwt(db, user)
	return *user.Access
}

// ----------------------------------------------------------------------------

// PROFILE FIXTURES
func CreateCity(db *gorm.DB) models.City {
	country := models.Country{Name: "Nigeria", Code: "NG"}
	db.FirstOrCreate(&country, models.Country{Name: country.Name})

	region := models.Region{Name: "Lagos", CountryId: country.ID}
	db.FirstOrCreate(&region, region)

	city := models.City{Name: "Lekki", CountryId: country.ID, RegionId: &region.ID}
	db.FirstOrCreate(&city, city)
	return city
}

func CreateFriend(db *gorm.DB, status choices.FriendStatusChoice) models.Friend {
	DropAndCreateSingleTable(db, models.Friend{})
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	friend := models.Friend{RequesterID: verifiedUser.ID, RequesteeID: anotherVerifiedUser.ID, Status: status}
	db.FirstOrCreate(&friend, friend)
	friend.Requester = verifiedUser
	friend.Requestee = anotherVerifiedUser
	return friend
}

func CreateNotification(db *gorm.DB) models.Notification {
	user := CreateTestVerifiedUser(db)
	text := "A new update is coming!"
	notification := notificationManager.Create(db, nil, "ADMIN", []models.User{user}, nil, nil, nil, &text)
	return notification
}

// ----------------------------------------------------------------------------

// CHAT FIXTURES
func CreateChat(db *gorm.DB) models.Chat {
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	chat := chatManager.GetDMChat(db, verifiedUser, anotherVerifiedUser)
	if chat.ID == nil {
		chat = chatManager.Create(db, verifiedUser, "DM", []models.User{anotherVerifiedUser})
	} else {
		// Set useful related data
		chat.OwnerObj = verifiedUser
	}
	chat.UserObjs = []models.User{anotherVerifiedUser}
	return chat
}

func CreateGroupChat(db *gorm.DB) models.Chat {
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	DropAndCreateSingleTable(db, models.Chat{})
	dataToCreate := schemas.GroupChatCreateSchema{Name: "My New Group"}
	chat := chatManager.CreateGroup(db, verifiedUser, []models.User{anotherVerifiedUser}, dataToCreate)
	return chat
}

func CreateMessage(db *gorm.DB) models.Message {
	DropAndCreateSingleTable(db, models.Message{})
	chat := CreateChat(db)
	text := "Hello Boss"
	message := messageManager.Create(db, chat.OwnerObj, chat, &text, nil)
	return message
}

// ----------------------------------------------------------------------------

// FEED FIXTURES
func CreatePost(db *gorm.DB) models.Post {
	author := CreateTestVerifiedUser(db)
	post := postManager.Create(db, author, schemas.PostInputSchema{Text: "This is a nice new platform."})
	return post
}

func CreateReaction(db *gorm.DB) models.Reaction {
	post := CreatePost(db)
	reaction := reactionManager.Create(db, post.AuthorObj, choices.FTPOST, &post, nil, nil, choices.RLIKE)
	return reaction
}

func CreateComment(db *gorm.DB) models.Comment {
	post := CreatePost(db)
	comment := commentManager.Create(db, post.AuthorObj, post, "Just a comment")
	return comment
}

func CreateReply(db *gorm.DB) models.Reply {
	comment := CreateComment(db)
	reply := replyManager.Create(db, comment.AuthorObj, comment, "Simple reply")
	return reply
}

// ----------------------------------------------------------------------------

// Utils
func GetUserMap(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"name":     user.FullName(),
		"username": user.Username,
		"avatar":   user.GetAvatarUrl(),
	}
}

func ConvertDateTime(timeObj time.Time) string {
	roundedTime := timeObj.Round(time.Microsecond)
	formatted := roundedTime.Format("2006-01-02T15:04:05")

	// Get the microsecond part and round it
	microseconds := roundedTime.Nanosecond() / 1000

	// Append the rounded microsecond part to the formatted string
	formatted = fmt.Sprintf("%s.%06d", formatted, microseconds)
	formatted = strings.TrimRight(formatted, "0")
	// Append the timezone information
	formatted = fmt.Sprintf("%s%s", formatted, roundedTime.Format("-07:00"))

	return formatted
}

// ----------------------------------------------------------------------------
