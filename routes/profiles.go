package routes

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"gorm.io/gorm/clause"
)

// @Summary Retrieve cities based on query params
// @Description This endpoint retrieves the first 10 cities that matches the query params
// @Tags Profiles
// @Param name query string false "City name"
// @Success 200 {object} schemas.CitiesResponseSchema
// @Router /profiles/cities [get]
func (endpoint Endpoint) RetrieveCities(c *fiber.Ctx) error {
	db := endpoint.DB
	message := "Cities Fetched"
	name := c.Query("name")

	// Define a regular expression to match non-word characters (excluding spaces).
	re := regexp.MustCompile(`[^\w\s]`)
	// Use the regular expression to replace matching substrings with an empty string.
	name = re.ReplaceAllString(name, "")
	cities := []models.City{}
	db.Preload(clause.Associations).Where("name ILIKE ?", "%"+name+"%").Find(&cities)

	if len(cities) == 0 {
		message = "No match found"
	}
	response := schemas.CitiesResponseSchema{
		ResponseSchema: SuccessResponse(message),
		Data:           cities,
	}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Retrieve Users
// @Description This endpoint retrieves a paginated list of users
// @Tags Profiles
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.ProfilesResponseSchema
// @Router /profiles [get]
// @Security BearerAuth
func (endpoint Endpoint) RetrieveUsers(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*models.User)

	users := []models.User{}
	query := db.Preload(clause.Associations)
	if user != nil {
		query.Not(models.User{BaseModel: models.BaseModel{ID: user.ID}})
	}
	query.Find(&users)
	// Paginate, Convert type and return Users
	paginatedData, paginatedUsers, err := PaginateQueryset(users, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	users = paginatedUsers.([]models.User)

	response := schemas.ProfilesResponseSchema{
		ResponseSchema: SuccessResponse("Users fetched"),
		Data: schemas.ProfilesResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       users,
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Retrieve User Profile
// @Description This endpoint retrieves a user profile
// @Tags Profiles
// @Param username path string true "Username of user"
// @Success 200 {object} schemas.ProfileResponseSchema
// @Router /profiles/profile/{username} [get]
func (endpoint Endpoint) RetrieveUserProfile(c *fiber.Ctx) error {
	db := endpoint.DB
	username := c.Params("username")

	user := models.User{}
	db.Preload(clause.Associations).Take(&user, models.User{Username: username})
	if user.ID == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No user with that username"))
	}

	// Return User
	response := schemas.ProfileResponseSchema{
		ResponseSchema: SuccessResponse("User details fetched"),
		Data:           user.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Profile
// @Description This endpoint updates a user profile
// @Tags Profiles
// @Param profile body schemas.ProfileUpdateSchema true "Profile object"
// @Success 200 {object} schemas.ProfileResponseSchema
// @Router /profiles/profile [patch]
// @Security BearerAuth
func (endpoint Endpoint) UpdateProfile(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*models.User)

	data := schemas.ProfileUpdateSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate City Value
	cityID := data.CityID
	if cityID != nil {
		city := models.City{}
		db.Take(&city, models.City{BaseModel: models.BaseModel{ID: *cityID}})
		if city.ID == nil {
			data := map[string]string{
				"city_id": "No city with that ID",
			}
			return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
		}
		user.CityObj = &city
	}

	// Create OR Update File
	fileType := data.FileType
	if fileType != nil {
		file := models.File{ResourceType: *fileType}.UpdateOrCreate(db, user.AvatarId)
		user.AvatarObj = &file
	}

	// Set values & save
	user = data.SetValues(user)
	db.Save(&user)

	// Return repsonse
	updatedData := schemas.ProfileUpdateResponseDataSchema{
		User: *user,
	}.Init(fileType)
	response := schemas.ProfileUpdateResponseSchema{
		ResponseSchema: SuccessResponse("User updated"),
		Data:           updatedData,
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete User's Account
// @Description This endpoint deletes a particular user's account (irreversible)
// @Tags Profiles
// @Param password body schemas.DeleteUserSchema true "Password"
// @Success 200 {object} schemas.ResponseSchema
// @Router /profiles/profile [post]
// @Security BearerAuth
func (endpoint Endpoint) DeleteUser(c *fiber.Ctx) error {
	db := endpoint.DB
	user := c.Locals("user").(*models.User)

	data := schemas.DeleteUserSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Check if password is valid
	if !utils.CheckPasswordHash(data.Password, user.Password) {
		data := map[string]string{
			"password": "Incorrect password",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Delete User
	db.Delete(&user)
	return c.Status(200).JSON(SuccessResponse("User deleted"))
}

// // @Summary Retrieve Friends
// // @Description This endpoint retrieves friends of a user
// // @Tags Profiles
// // @Param page query int false "Current Page" default(1)
// // @Success 200 {object} schemas.ProfilesResponseSchema
// // @Router /profiles/friends [get]
// // @Security BearerAuth
// func (endpoint Endpoint) RetrieveFriends(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	friends := friendManager.GetFriends(db, user)

// 	// Paginate, Convert type and return Friends
// 	paginatedData, paginatedFriends, err := PaginateQueryset(friends, c, 20)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	convertedFriends := utils.ConvertStructData(paginatedFriends, []schemas.ProfileSchema{}).(*[]schemas.ProfileSchema)
// 	response := schemas.ProfilesResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Friends fetched"}.Init(),
// 		Data: schemas.ProfilesResponseDataSchema{
// 			PaginatedResponseDataSchema: *paginatedData,
// 			Items:                       *convertedFriends,
// 		}.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Retrieve Friend Requests
// // @Description This endpoint retrieves friend requests of a user
// // @Tags Profiles
// // @Param page query int false "Current Page" default(1)
// // @Success 200 {object} schemas.ProfilesResponseSchema
// // @Router /profiles/friends/requests [get]
// // @Security BearerAuth
// func (endpoint Endpoint) RetrieveFriendRequests(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	friendsRequests := friendManager.GetFriendRequests(db, user)

// 	// Paginate, Convert type and return Friends Requests
// 	paginatedData, paginatedFriendRequests, err := PaginateQueryset(friendsRequests, c, 20)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	convertedFriendRequests := utils.ConvertStructData(paginatedFriendRequests, []schemas.ProfileSchema{}).(*[]schemas.ProfileSchema)
// 	response := schemas.ProfilesResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Friend Requests fetched"}.Init(),
// 		Data: schemas.ProfilesResponseDataSchema{
// 			PaginatedResponseDataSchema: *paginatedData,
// 			Items:                       *convertedFriendRequests,
// 		}.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Send Or Delete Friend Request
// // @Description This endpoint sends or delete friend requests
// // @Tags Profiles
// // @Param friend_request body schemas.SendFriendRequestSchema true "Friend Request object"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /profiles/friends/requests [post]
// // @Security BearerAuth
// func (endpoint Endpoint) SendOrDeleteFriendRequest(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	friendRequestData := schemas.SendFriendRequestSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &friendRequestData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(friendRequestData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	requestee, friend, errData := friendManager.GetRequesteeAndFriendObj(db, user, friendRequestData.Username)
// 	if errData != nil {
// 		return c.Status(404).JSON(errData)
// 	}
// 	message := "Friend Request sent"
// 	statusCode := 201

// 	if friend != nil {
// 		statusCode = 200
// 		message = "Friend Request removed"
// 		if friend.Status == "ACCEPTED" {
// 			message = "This user is already your friend"
// 		} else if user.ID != friend.RequesterID {
// 			return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "The user already sent you a friend request!"))
// 		} else {
// 			// Delete friend successfully
// 			db.Friend.DeleteOne(friend).Exec(managers.Ctx)
// 		}

// 	} else {
// 		// Create Friend Object
// 		friendManager.Create(db, user, requestee, "PENDING")
// 	}

// 	response := schemas.ResponseSchema{Message: message}.Init()
// 	return c.Status(statusCode).JSON(response)
// }

// // @Summary Accept Or Reject a Friend Request
// // @Description This endpoint accepts or reject a friend request
// // @Tags Profiles
// // @Param friend_request body schemas.AcceptFriendRequestSchema true "Friend Request object"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /profiles/friends/requests [put]
// // @Security BearerAuth
// func (endpoint Endpoint) AcceptOrRejectFriendRequest(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	friendRequestData := schemas.AcceptFriendRequestSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &friendRequestData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(friendRequestData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	_, friend, errData := friendManager.GetRequesteeAndFriendObj(db, user, friendRequestData.Username, "PENDING")
// 	if errData != nil {
// 		return c.Status(404).JSON(errData)
// 	}
// 	if friend == nil {
// 		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No friend request exist between you and that user"))
// 	}
// 	if friend.RequesterID == user.ID {
// 		return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "You cannot accept or reject a friend request you sent"))
// 	}
// 	// Update or delete friend request based on status
// 	message := "Accepted"
// 	if friendRequestData.Accepted {
// 		// Update Friend Request
// 		friend.Update().SetStatus("ACCEPTED").Save(managers.Ctx)
// 	} else {
// 		// Delete Friend Request
// 		message = "Rejected"
// 		db.Friend.DeleteOne(friend).Exec(managers.Ctx)
// 	}
// 	response := schemas.ResponseSchema{Message: fmt.Sprintf("Friend Request %s", message)}.Init()
// 	return c.Status(200).JSON(response)
// }

// var notificationManager = managers.NotificationManager{}
// // @Summary Retrieve User Notifications
// // @Description This endpoint retrieves a paginated list of auth user's notifications. Use post, comment, reply slug to navigate to the post, comment or reply.
// // @Tags Profiles
// // @Param page query int false "Current Page" default(1)
// // @Success 200 {object} schemas.NotificationsResponseSchema
// // @Router /profiles/notifications [get]
// // @Security BearerAuth
// func (endpoint Endpoint) RetrieveUserNotifications(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	notifications := notificationManager.GetQueryset(db, user.ID)

// 	// Paginate, Convert type and return notifications
// 	paginatedData, paginatedNotifications, err := PaginateQueryset(notifications, c)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	convertedNotifications := utils.ConvertStructData(paginatedNotifications, []schemas.NotificationSchema{}).(*[]schemas.NotificationSchema)
// 	response := schemas.NotificationsResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Notifications fetched"}.Init(),
// 		Data: schemas.NotificationsResponseDataSchema{
// 			PaginatedResponseDataSchema: *paginatedData,
// 			Items:                       *convertedNotifications,
// 		}.Init(user.ID),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Read Notifications
// // @Description This endpoint reads a notification
// // @Tags Profiles
// // @Param read_data body schemas.ReadNotificationSchema true "Read Notification Data"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /profiles/notifications [post]
// // @Security BearerAuth
// func (endpoint Endpoint) ReadNotification(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	readNotificationData := schemas.ReadNotificationSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &readNotificationData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(readNotificationData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	notificationID := readNotificationData.ID
// 	markAllAsRead := readNotificationData.MarkAllAsRead

// 	respMessage := "Notifications read"
// 	if markAllAsRead {
//         // Mark all notifications as read
// 		notificationManager.MarkAsRead(db, user.ID)
// 	} else if notificationID != nil {
//         // Mark single notification as read
// 		err := notificationManager.ReadOne(db, user.ID, *notificationID)
// 		if err != nil {
// 			return c.Status(404).JSON(err)
// 		}
// 		respMessage = "Notification read"
// 	}
// 	response := schemas.ResponseSchema{Message: respMessage}.Init()
// 	return c.Status(200).JSON(response)
// }
