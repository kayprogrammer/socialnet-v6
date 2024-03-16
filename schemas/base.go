package schemas

import (
	"github.com/kayprogrammer/socialnet-v6/models"
)

type ResponseSchema struct {
	Status		string		`json:"status" example:"success"`
	Message		string		`json:"message" example:"Data fetched/created/updated/deleted"`
}

type PaginatedResponseDataSchema struct {
	PerPage     uint `json:"per_page" example:"100"`
	CurrentPage uint `json:"current_page" example:"1"`
	LastPage    uint `json:"last_page" example:"100"`
}

type UserDataSchema struct {
	Name     string  `json:"name" example:"John Doe"`
	Username string  `json:"username" example:"john-doe"`
	Avatar   *string `json:"avatar" example:"https://img.url"`
}

func (user UserDataSchema) Init(userObj *models.User) UserDataSchema {
	user.Name = userObj.FullName()
	user.Username = userObj.Username
	user.Avatar = userObj.GetAvatarUrl()
	return user
}
