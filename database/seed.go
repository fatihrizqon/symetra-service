package database

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedDefaultUser(db *gorm.DB) {
	var count int64

	db.Model(&entity.User{}).Count(&count)

	if count > 0 {
		return
	}

	user := entity.User{
		Name:     "Administrator",
		Username: "administrator",
		Email:    "admin@example.com",
		Password: hashPassword("password"),
	}

	db.Create(&user)
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
