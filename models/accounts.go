package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"gorm.io/gorm"
)

type Country struct {
	BaseModel
	Name string `json:"name" gorm:"type: varchar(255);not null" example:"Nigeria"`
	Code string `json:"code" gorm:"type: varchar(255);not null" example:"NG"`
}

type Region struct {
	BaseModel
	Name      string    `json:"name" gorm:"type: varchar(255);not null" example:"Lagos"`
	CountryId uuid.UUID `json:"country_id" gorm:"not null"`
	Country   Country   `gorm:"foreignKey:CountryId;constraint:OnDelete:CASCADE"`
}

type City struct {
	BaseModel
	Name      string     `json:"name" gorm:"type: varchar(255);not null" example:"Lekki"`
	RegionId  *uuid.UUID `json:"region_id"`
	Region    *Region    `gorm:"foreignKey:RegionId;constraint:OnDelete:SET NULL"`
	CountryId uuid.UUID  `json:"country_id" gorm:"not null"`
	Country   Country    `gorm:"foreignKey:CountryId;constraint:OnDelete:CASCADE"`
}

type User struct {
	BaseModel
	FirstName       string     `json:"first_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"John"`
	LastName        string     `json:"last_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"Doe"`
	Username        string     `json:"username" gorm:"type: varchar(1000);not null;unique;" validate:"required,max=255" example:"john-doe"`
	Email           string     `json:"email" gorm:"not null;unique;" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password        string     `json:"password" gorm:"not null" validate:"required,min=8,max=50" example:"strongpassword"`
	IsEmailVerified bool       `json:"is_email_verified" gorm:"default:false" swaggerignore:"true"`
	IsSuperuser     bool       `json:"is_superuser" gorm:"default:false" swaggerignore:"true"`
	IsStaff         bool       `json:"is_staff" gorm:"default:false" swaggerignore:"true"`
	TermsAgreement  bool       `json:"terms_agreement" gorm:"default:false" validate:"eq=true"`
	AvatarId        *uuid.UUID `json:"avatar_id" gorm:"null" swagger:"ignore" swaggerignore:"true"`
	Avatar          *File      `gorm:"foreignKey:AvatarId;constraint:OnDelete:SET NULL;null;" swaggerignore:"true"`
	Access          *string    `gorm:"type:varchar(1000);null;" json:"access"`
	Refresh         *string    `gorm:"type:varchar(1000);null;" json:"refresh"`
	Bio             *string    `gorm:"type:varchar(1000);null;" json:"bio"`
	Dob             *time.Time `gorm:"null;" json:"dob"`
	CityId          *uuid.UUID `json:"city_id" gorm:"null"`
	City            *City      `gorm:"foreignKey:CityId;constraint:OnDelete:SET NULL"`
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Hash password
	user.Password = utils.HashPassword(user.Password)

	// Create username
	user.Username = user.GenerateUsername(tx)
	return
}

func (user *User) GenerateUsername(tx *gorm.DB) (string) {
	uniqueUsername := slug.Make(user.FirstName + " " + user.LastName)
	userName := user.Username
	if userName != "" {
		uniqueUsername = userName
	}

	existingUser := User{Username: uniqueUsername}
	tx.Take(&existingUser, existingUser)
	if existingUser.ID != uuid.Nil { // username is already taken
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
