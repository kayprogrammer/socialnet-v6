package tests

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kayprogrammer/socialnet-v6/managers"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/routes"
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
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	friend := models.Friend{RequesterID: verifiedUser.ID, RequesteeID: anotherVerifiedUser.ID}
	friend := friendManager.Create(db, verifiedUser, anotherVerifiedUser, status)
	return friend
}

func CreateNotification(db *gorm.DB) *ent.Notification {
	user := CreateTestVerifiedUser(db)
	text := "A new update is coming!"
	notification := notificationManager.Create(db, nil, "ADMIN", []uuid.UUID{user.ID}, nil, nil, nil, &text)
	return notification
}

// ----------------------------------------------------------------------------

// CHAT FIXTURES
func CreateChat(db *gorm.DB) *ent.Chat {
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	chat := chatManager.GetDMChat(db, verifiedUser, anotherVerifiedUser)
	if chat == nil {
		chat = chatManager.Create(db, verifiedUser, "DM", []*ent.User{anotherVerifiedUser})
	} else {
		// Set useful related data
		chat.Edges.Owner = verifiedUser
	}
	chat.Edges.Users = []*ent.User{anotherVerifiedUser}
	return chat
}

func CreateGroupChat(db *gorm.DB) *ent.Chat {
	verifiedUser := CreateTestVerifiedUser(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)
	chatManager.DropData(db)
	dataToCreate := schemas.GroupChatCreateSchema{Name: "My New Group"}
	chat := chatManager.CreateGroup(db, verifiedUser, []*ent.User{anotherVerifiedUser}, dataToCreate)
	chat.Edges.Users = []*ent.User{anotherVerifiedUser}
	return chat
}

func CreateMessage(db *gorm.DB) *ent.Message {
	messageManager.DropData(db)
	chat := CreateChat(db)
	text := "Hello Boss"
	message := messageManager.Create(db, chat.Edges.Owner, chat, &text, nil)
	return message
}

// ----------------------------------------------------------------------------

// FEED FIXTURES
func CreatePost(db *gorm.DB) *ent.Post {
	author := CreateTestVerifiedUser(db)
	post := postManager.Create(db, author, schemas.PostInputSchema{Text: "This is a nice new platform."})
	return post
}

func CreateReaction(db *gorm.DB) *ent.Reaction {
	post := CreatePost(db)
	reaction := reactionManager.Create(db, post.AuthorID, "POST", post.ID, "LIKE")
	reaction.Edges.Post = post
	reaction.Edges.User = post.Edges.Author
	return reaction
}

func CreateComment(db *gorm.DB) *ent.Comment {
	post := CreatePost(db)
	comment := commentManager.Create(db, post.Edges.Author, post.ID, "Just a comment")
	comment.Edges.Post = post
	return comment
}

func CreateReply(db *gorm.DB) *ent.Reply {
	comment := CreateComment(db)
	reply := replyManager.Create(db, comment.Edges.Author, comment.ID, "Simple reply")
	reply.Edges.Comment = comment
	return reply
}

// ----------------------------------------------------------------------------

// Utils
func GetUserMap(user *ent.User) map[string]interface{} {
	return map[string]interface{}{
		"name":     schemas.FullName(user),
		"username": user.Username,
		"avatar":   nil,
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
