package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/schemas"
)

// @Summary Retrieve site details
// @Description This endpoint retrieves few details of the site/application.
// @Tags General
// @Success 200 {object} schemas.SiteDetailResponseSchema
// @Router /general/site-detail [get]
func (ep Endpoint) GetSiteDetails(c *fiber.Ctx) error {
	db := ep.DB
	var sitedetail models.SiteDetail

	db.FirstOrCreate(&sitedetail, sitedetail)
	responseSiteDetail := schemas.SiteDetailResponseSchema{
		ResponseSchema: SuccessResponse("Site Details Fetched!"),
		Data:           sitedetail,
	}
	return c.Status(200).JSON(responseSiteDetail)
}
