package repository

import (
	"time"

	"gorm.io/gorm"
)

// AccountLedgerRow is a raw row from the aggregated ledger query.
type AccountLedgerRow struct {
	AccountCode  string
	AccountName  string
	GroupName    string // e.g. "Assets", "Liabilities", "Equity", "Revenue", "Expense"
	SubgroupName string
	TotalDebit   float64
	TotalCredit  float64
	PostedCount  int
}

type IReportRepository interface {
	GetLedger(start, end time.Time) ([]AccountLedgerRow, error)
	GetLedgerUpTo(end time.Time) ([]AccountLedgerRow, error)
	GetPostedCount(start, end time.Time) (int, error)
}

type ReportRepository struct {
	Db *gorm.DB
}

func NewReportRepository(db *gorm.DB) IReportRepository {
	return &ReportRepository{Db: db}
}

// GetLedger aggregates all posted journal lines within a date range.
func (r *ReportRepository) GetLedger(start, end time.Time) ([]AccountLedgerRow, error) {
	var rows []AccountLedgerRow

	err := r.Db.Raw(`
		SELECT
			ca.code                      AS account_code,
			ca.name                      AS account_name,
			cg.name                      AS group_name,
			cs.name                      AS subgroup_name,
			COALESCE(SUM(jl.debit),  0)  AS total_debit,
			COALESCE(SUM(jl.credit), 0)  AS total_credit,
			COUNT(DISTINCT je.id)        AS posted_count
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		INNER JOIN coa_subgroups cs   ON cs.id = ca.subgroup_id
		INNER JOIN coa_groups cg      ON cg.id = cs.group_id
		WHERE je.status = 'posted'
		  AND je.date >= ?
		  AND je.date <= ?
		GROUP BY ca.code, ca.name, cg.name, cs.name
		ORDER BY ca.code
	`, start, end).Scan(&rows).Error

	return rows, err
}

// GetLedgerUpTo aggregates all posted journal lines from beginning of time up to end (inclusive).
// Used for Balance Sheet (cumulative) and opening balances.
func (r *ReportRepository) GetLedgerUpTo(end time.Time) ([]AccountLedgerRow, error) {
	var rows []AccountLedgerRow

	err := r.Db.Raw(`
		SELECT
			ca.code                      AS account_code,
			ca.name                      AS account_name,
			cg.name                      AS group_name,
			cs.name                      AS subgroup_name,
			COALESCE(SUM(jl.debit),  0)  AS total_debit,
			COALESCE(SUM(jl.credit), 0)  AS total_credit,
			COUNT(DISTINCT je.id)        AS posted_count
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		INNER JOIN coa_subgroups cs   ON cs.id = ca.subgroup_id
		INNER JOIN coa_groups cg      ON cg.id = cs.group_id
		WHERE je.status = 'posted'
		  AND je.date <= ?
		GROUP BY ca.code, ca.name, cg.name, cs.name
		ORDER BY ca.code
	`, end).Scan(&rows).Error

	return rows, err
}

// GetPostedCount returns number of posted entries in range.
func (r *ReportRepository) GetPostedCount(start, end time.Time) (int, error) {
	var count int64
	err := r.Db.Raw(`
		SELECT COUNT(*) FROM journal_entries
		WHERE status = 'posted' AND date >= ? AND date <= ?
	`, start, end).Scan(&count).Error
	return int(count), err
}
