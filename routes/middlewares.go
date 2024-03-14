package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"gorm.io/gorm"
)

func getUser(c *fiber.Ctx, token string, db *gorm.DB) (*models.User, *string) {
	if len(token) < 8 {
		err := "Auth Token is Invalid or Expired!"
		return nil, &err
	}
	user, err := DecodeAccessToken(token[7:], db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ep Endpoint) AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB

	if len(token) < 1 {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
	}
	user, err := getUser(c, token, db)
	if err != nil {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, "Access token is invalid or expired!"))
	}
	c.Locals("user", user)
	return c.Next()
}

func ParseUUID(input string) *uuid.UUID {
    uuidVal, err := uuid.Parse(input)
	if err != nil {
		return nil
	}
    return &uuidVal
}