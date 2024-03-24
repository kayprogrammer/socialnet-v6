package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
)

func SuccessResponse(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func ValidateReactionFocus(focus string) *utils.ErrorResponse {
	switch focus {
	case "POST", "COMMENT", "REPLY":
		return nil
	}
	err := utils.RequestErr(utils.ERR_INVALID_VALUE, "Invalid 'focus' value")
	return &err
}