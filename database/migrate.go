package database

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entity.User{},
		&entity.Session{},
		&entity.Credential{},
		&entity.COAGroup{},
		&entity.COASubGroup{},
	)
	seedDefaultUser(db)
}
