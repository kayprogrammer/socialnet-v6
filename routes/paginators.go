package routes

import (
	"math"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
)

func PaginateQueryset(queryset interface{}, fiberCtx *fiber.Ctx, opts ...int) (*schemas.PaginatedResponseDataSchema, any, *utils.ErrorResponse) {
	currentPage := fiberCtx.QueryInt("page", 1)
	var perPage int

	if currentPage < 1 {
		errData := utils.RequestErr(utils.ERR_INVALID_PAGE, "Invalid Page")
		return nil, nil, &errData
	}

	// Check if page size is provided as an argument
	if len(opts) > 0 {
		perPage = opts[0]
	} else {
		// Default page size if not provided
		perPage = 50
	}
	querysetValue := reflect.ValueOf(queryset)
	itemsCount := querysetValue.Len()
	lastPage := math.Ceil(float64(itemsCount) / float64(perPage))
	if lastPage == 0 {
		lastPage = 1
	}
	if currentPage > int(lastPage) {
		errData := utils.RequestErr(utils.ERR_INVALID_PAGE, "Page number is out of range")
		return nil, nil, &errData
	}

	startIndex := (currentPage - 1) * perPage
	endIndex := startIndex + perPage

	// Ensure startIndex is within bounds
	if startIndex < 0 {
		startIndex = 0
	}

	// Ensure endIndex is within bounds
	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	paginatorData := schemas.PaginatedResponseDataSchema{
		PerPage:     uint(perPage),
		CurrentPage: uint(currentPage),
		LastPage:    uint(lastPage),
	}
	paginatedItems := querysetValue.Slice(startIndex, endIndex).Interface()
	return &paginatorData, paginatedItems, nil
}
