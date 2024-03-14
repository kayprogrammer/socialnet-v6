package models

type SiteDetail struct {
	BaseModel
	Name    string `json:"name" gorm:"default:SocialNet;type:varchar(50);not null"`
	Email   string `json:"email" gorm:"default:kayprogrammer1@gmail.com;not null" example:"kayprogrammer1@gmail.com"`
	Phone   string `json:"phone" gorm:"default:+2348133831036;type:varchar(20);not null" example:"+2348133831036"`
	Address string `json:"address" gorm:"default:234, Lagos, Nigeria;not null" example:"234, Lagos, Nigeria"`
	Fb      string `json:"fb" gorm:"default:https://facebook.com;not null" example:"https://facebook.com"`
	Tw      string `json:"tw" gorm:"default:https://twitter.com;not null" example:"https://twitter.com"`
	Wh      string `json:"wh" gorm:"default:https://wa.me/2348133831036;not null" example:"https://wa.me/2348133831036"`
	Ig      string `json:"ig" gorm:"default:https://instagram.com;not null" example:"https://instagram.com"`
}

