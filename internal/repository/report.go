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
	GetGeneralLedger(start, end time.Time, coaID string) ([]GeneralLedgerRow, error)
	GetOpeningBalance(asOf time.Time, coaID string) (float64, error)
	GetJournalBook(start, end time.Time) ([]JournalBookRow, error)
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

// GeneralLedgerRow is a single transaction line for the general ledger report.
type GeneralLedgerRow struct {
	AccountCode    string
	AccountName    string
	GroupName      string
	SubgroupName   string
	Date           string
	JournalNumber  string
	Description    string
	Debit          float64
	Credit         float64
}

// JournalBookRow is a single journal line for the journal book report.
type JournalBookRow struct {
	Date          string
	JournalNumber string
	JournalType   string
	Description   string
	AccountCode   string
	AccountName   string
	Debit         float64
	Credit        float64
	TotalDebit    float64
	TotalCredit   float64
}

// GetGeneralLedger fetches all posted transaction lines per account for the period.
// Pass coaID = "" to fetch all accounts, or a specific UUID to filter one account.
func (r *ReportRepository) GetGeneralLedger(start, end time.Time, coaID string) ([]GeneralLedgerRow, error) {
	var rows []GeneralLedgerRow

	query := `
		SELECT
			ca.code                             AS account_code,
			ca.name                             AS account_name,
			cg.name                             AS group_name,
			cs.name                             AS subgroup_name,
			TO_CHAR(je.date, 'YYYY-MM-DD')      AS date,
			je.journal_number                   AS journal_number,
			je.description                      AS description,
			COALESCE(jl.debit,  0)              AS debit,
			COALESCE(jl.credit, 0)              AS credit
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		INNER JOIN coa_subgroups cs   ON cs.id = ca.subgroup_id
		INNER JOIN coa_groups cg      ON cg.id = cs.group_id
		WHERE je.status = 'posted'
		  AND je.date >= ?
		  AND je.date <= ?`

	args := []interface{}{start, end}
	if coaID != "" {
		query += " AND ca.id = ?"
		args = append(args, coaID)
	}
	query += " ORDER BY ca.code, je.date, je.journal_number"

	err := r.Db.Raw(query, args...).Scan(&rows).Error
	return rows, err
}

// GetOpeningBalance returns the net opening balance of an account (or all accounts
// if coaID = "") cumulatively up to (but not including) the start date.
// Normal balance sign: assets/expenses → debit−credit, liabilities/equity/revenue → credit−debit.
// We return raw (debit − credit) and let the service apply the sign per account type.
func (r *ReportRepository) GetOpeningBalance(asOf time.Time, coaID string) (float64, error) {
	// asOf is exclusive — we want everything BEFORE the period start
	type result struct{ Balance float64 }
	var res result

	query := `
		SELECT COALESCE(SUM(jl.debit) - SUM(jl.credit), 0) AS balance
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		WHERE je.status = 'posted'
		  AND je.date < ?`

	args := []interface{}{asOf}
	if coaID != "" {
		query += " AND ca.id = ?"
		args = append(args, coaID)
	}

	err := r.Db.Raw(query, args...).Scan(&res).Error
	return res.Balance, err
}

// GetJournalBook fetches all posted journal entries with their lines for the period.
func (r *ReportRepository) GetJournalBook(start, end time.Time) ([]JournalBookRow, error) {
	var rows []JournalBookRow

	err := r.Db.Raw(`
		SELECT
			TO_CHAR(je.date, 'YYYY-MM-DD')  AS date,
			je.journal_number               AS journal_number,
			je.type                         AS journal_type,
			je.description                  AS description,
			ca.code                         AS account_code,
			ca.name                         AS account_name,
			COALESCE(jl.debit,  0)          AS debit,
			COALESCE(jl.credit, 0)          AS credit,
			je.total_debit                  AS total_debit,
			je.total_credit                 AS total_credit
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		WHERE je.status = 'posted'
		  AND je.date >= ?
		  AND je.date <= ?
		ORDER BY je.date, je.journal_number, jl.id
	`, start, end).Scan(&rows).Error

	return rows, err
}
