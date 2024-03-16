package initials

import (
	"log"

	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/models"
	"gorm.io/gorm"
)

func createSuperUser(cfg config.Config, db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Email:           cfg.FirstSuperuserEmail,
		Password:        cfg.FirstSuperUserPassword,
		IsSuperuser:     true,
		IsStaff:         true,
		IsEmailVerified: true,
		TermsAgreement: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func createClient(cfg config.Config, db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Email:           cfg.FirstClientEmail,
		Password:        cfg.FirstClientPassword,
		IsEmailVerified: true,
		TermsAgreement: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func createCity(cfg config.Config, db *gorm.DB) models.City {
	country := models.Country{
		Name: "Nigeria",
		Code: "NG",
	}
	db.FirstOrCreate(&country, country)

	region := models.Region{
		Name: "Lagos",
		CountryId: country.ID,
	}
	db.FirstOrCreate(&region, models.Region{Name: "Lagos"})

	city := models.City{
		Name: "Lekki",
		RegionId: &region.ID,
		CountryId: country.ID,
	}
	db.FirstOrCreate(&city, models.City{Name: "Lekki"})
	return city
}

func CreateInitialData(cfg config.Config, db *gorm.DB) {
	log.Println("Creating Initial Data....")
	createSuperUser(cfg, db)
	createClient(cfg, db)
	createCity(cfg, db)
	log.Println("Initial Data Created....")
}