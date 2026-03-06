package database

import (
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		// ── Auth & Users ──────────────────────────────────────────────────────
		&entity.User{},
		&entity.Session{},
		&entity.Credential{},

		// ── Chart of Accounts hierarchy ───────────────────────────────────────
		&entity.COAGroup{},
		&entity.COASubGroup{},
		&entity.COA{},

		// ── Fiscal ────────────────────────────────────────────────────────────
		// FiscalYear must come before FiscalPeriod (FK dependency)
		// FiscalPeriodLog must come after FiscalPeriod
		&entity.FiscalYear{},
		&entity.FiscalPeriod{},
		&entity.FiscalPeriodLog{},

		// ── Transactions ──────────────────────────────────────────────────────
		&entity.JournalEntry{},
		&entity.JournalLine{},
	)

	seedDefaultData(db)
}

func seedDefaultData(db *gorm.DB) {
	SeedData(db)
}
