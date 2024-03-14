package routes

import "github.com/kayprogrammer/socialnet-v6/schemas"

func SuccessResponse(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}