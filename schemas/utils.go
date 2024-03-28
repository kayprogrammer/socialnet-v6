package schemas

import "github.com/kayprogrammer/socialnet-v6/models"

func ConvertUsers(users []models.User) []UserDataSchema {
	convertedUsers := []UserDataSchema{}
	for i := range users {
		user := UserDataSchema{}.Init(users[i])
		convertedUsers = append(convertedUsers, user)
	}
	return convertedUsers
}