package schemas

import "github.com/kayprogrammer/socialnet-v6/models"

type SiteDetailResponseSchema struct {
	ResponseSchema
	Data models.SiteDetail `json:"data"`
}
