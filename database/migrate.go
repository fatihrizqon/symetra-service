package database

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		// Auth & Users
		&entity.User{},
		&entity.Session{},
		&entity.Credential{},

		// Chart of Accounts hierarchy
		&entity.COAGroup{},
		&entity.COASubGroup{},
		&entity.COA{},

		// Master data
		&entity.Customer{},
		&entity.Vendor{},

		// Journal entries — must come after COA
		&entity.JournalEntry{},
		&entity.JournalLine{},
	)

	seedDefaultUser(db)
}
