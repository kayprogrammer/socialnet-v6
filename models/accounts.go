package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

type Country struct {
	BaseModel
	Name string `json:"name" gorm:"type: varchar(255);not null;unique" example:"Nigeria"`
	Code string `json:"code" gorm:"type: varchar(255);not null;unique" example:"NG"`
}

type Region struct {
	BaseModel
	Name      string    `json:"name" gorm:"type: varchar(255);not null" example:"Lagos"`
	CountryId uuid.UUID `json:"country_id" gorm:"not null"`
	Country   Country   `gorm:"foreignKey:CountryId;constraint:OnDelete:CASCADE"`
}

type City struct {
	BaseModel
	Name       string     `json:"name" gorm:"type: varchar(255);not null" example:"Lekki"`
	RegionId   *uuid.UUID `json:"-"`
	RegionObj  *Region    `json:"-" gorm:"foreignKey:RegionId;constraint:OnDelete:SET NULL"`
	Region     *string    `json:"region" gorm:"-" example:"Lagos"`
	CountryId  uuid.UUID  `json:"-" gorm:"not null"`
	CountryObj Country    `json:"-" gorm:"foreignKey:CountryId;constraint:OnDelete:CASCADE"`
	Country    string     `json:"country" gorm:"-" example:"Nigeria"`
}

func (city City) Init() City {
	// Set Related Data.
	region := city.RegionObj
	if region != nil {
		city.Region = &region.Name
	}
	city.Country = city.CountryObj.Name
	return city
}

type User struct {
	BaseModel
	FirstName       string     `json:"first_name" gorm:"type: varchar(255);not null" example:"John"`
	LastName        string     `json:"last_name" gorm:"type: varchar(255);not null" example:"Doe"`
	Username        string     `json:"username" gorm:"type: varchar(1000);not null;unique;" example:"john-doe"`
	Email           string     `json:"email" gorm:"not null;unique;" example:"johndoe@email.com"`
	Password        string     `json:"-" gorm:"not null"`
	IsEmailVerified bool       `json:"-" gorm:"default:false"`
	IsSuperuser     bool       `json:"-" gorm:"default:false"`
	IsStaff         bool       `json:"-" gorm:"default:false"`
	TermsAgreement  bool       `json:"-" gorm:"default:false"`
	AvatarId        *uuid.UUID `json:"-" gorm:"null"`
	AvatarObj       *File      `json:"-" gorm:"foreignKey:AvatarId;constraint:OnDelete:SET NULL;null;"`
	Avatar          *string    `json:"avatar"`
	Access          *string    `gorm:"type:varchar(1000);null;" json:"-"`
	Refresh         *string    `gorm:"type:varchar(1000);null;" json:"-"`
	Bio             *string    `gorm:"type:varchar(1000);null;" json:"bio"`
	Dob             *time.Time `gorm:"null;" json:"dob"`
	CityId          *uuid.UUID `json:"-" gorm:"null"`
	CityObj         *City      `json:"-" gorm:"foreignKey:CityId;constraint:OnDelete:SET NULL"`
	City            *string    `json:"city"`
}

func (user User) Init() User {
	user.ID = nil // Omit ID
	user.Avatar = user.GetAvatarUrl()
	user.City = user.GetCityName()
	return user
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user User) GetCityName() *string {
	city := user.CityObj
	if city != nil {
		return &city.Name
	}
	return nil
}

func (user User) GetAvatarUrl() *string {
	avatar := user.AvatarObj
	if avatar != nil {
		url := utils.GenerateFileUrl(avatar.ID.String(), "avatars", avatar.ResourceType)
		return &url
	}
	return nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Hash password
	user.Password = utils.HashPassword(user.Password)

	// Create username
	user.Username = user.GenerateUsername(tx)
	return
}

func (user *User) GenerateUsername(tx *gorm.DB) string {
	uniqueUsername := slug.Make(user.FirstName + " " + user.LastName)
	userName := user.Username
	if userName != "" {
		uniqueUsername = userName
	}

	existingUser := User{Username: uniqueUsername}
	tx.Take(&existingUser, existingUser)
	if existingUser.ID != nil { // username is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		user.Username = uniqueUsername + "-" + randomStr
		return user.GenerateUsername(tx)
	}
	return uniqueUsername
}

type Otp struct {
	BaseModel
	UserId uuid.UUID `json:"user_id" gorm:"unique"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Code   uint32    `json:"code"`
}

func (otp *Otp) BeforeSave(tx *gorm.DB) (err error) {
	code := uint32(utils.GetRandomInt(6))
	otp.Code = code
	return
}

func (obj Otp) CheckExpiration() bool {
	cfg := config.GetConfig()
	currentTime := time.Now().UTC()
	diff := int64(currentTime.Sub(obj.UpdatedAt).Seconds())
	emailExpirySecondsTimeout := cfg.EmailOtpExpireSeconds
	return diff > emailExpirySecondsTimeout
}
