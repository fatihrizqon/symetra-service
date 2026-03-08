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

		// ── Company (must precede all company-scoped tables) ──────────────────
		&entity.Company{},
		&entity.CompanyMember{},

		// ── Chart of Accounts hierarchy (now company-scoped) ──────────────────
		&entity.COAGroup{},
		&entity.COASubGroup{},
		&entity.COA{},

		// ── Transactional data (company-scoped) ───────────────────────────────
		&entity.JournalEntry{},
		&entity.JournalLine{},

		// ── Fiscal (company-scoped) ───────────────────────────────────────────
		&entity.FiscalYear{},
		&entity.FiscalPeriod{},
		&entity.FiscalPeriodLog{},

		// ── Parties (company-scoped) ──────────────────────────────────────────
		// Customer and Vendor handlers will be added in a subsequent phase.
		// Entities are migrated now so foreign keys resolve correctly.
		&entity.Customer{},
		&entity.Vendor{},

		// ── Company Configuration ─────────────────────────────────────────────
		&entity.CompanyConfiguration{},

		// ── Sales (company-scoped) ────────────────────────────────────────────
		&entity.Quotation{},
		&entity.QuotationItem{},
		&entity.Invoice{},
		&entity.InvoiceItem{},
	)

	seedDefaultData(db)
}

func seedDefaultData(db *gorm.DB) {
	SeedData(db)
}
