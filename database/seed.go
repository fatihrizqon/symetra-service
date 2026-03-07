package database

import (
	"log"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	log.Println("[SEED] Starting seed data...")

	seedAdminUser(db)

	log.Println("[SEED] Seed complete.")
}

func seedAdminUser(db *gorm.DB) {
	var user entity.User
	if db.Where("email = ?", "admin@symetra.id").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "admin", Name: "System Administrator",
			Email:        "admin@symetra.id",
			IsSuperadmin: true,
			Password:     string(hash), Status: 1,
		})
		log.Println("[SEED] Admin user created: admin@symetra.id / password")
	}
	if db.Where("email = ?", "operator@symetra.id").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "operator", Name: "System Operator",
			Email:    "operator@symetra.id",
			Password: string(hash), Status: 0,
		})
		log.Println("[SEED] operator user created: operator@symetra.id / password")
	}
}
