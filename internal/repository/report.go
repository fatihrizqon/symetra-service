package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountLedgerRow is a raw row from the aggregated ledger query.
type AccountLedgerRow struct {
	AccountCode  string
	AccountName  string
	GroupName    string
	SubgroupName string
	TotalDebit   float64
	TotalCredit  float64
	PostedCount  int
}

type IReportRepository interface {
	GetLedger(companyID uuid.UUID, start, end time.Time) ([]AccountLedgerRow, error)
	GetLedgerUpTo(companyID uuid.UUID, end time.Time) ([]AccountLedgerRow, error)
	GetPostedCount(companyID uuid.UUID, start, end time.Time) (int, error)
	GetGeneralLedger(companyID uuid.UUID, start, end time.Time, coaID string) ([]GeneralLedgerRow, error)
	GetOpeningBalance(companyID uuid.UUID, asOf time.Time, coaID string) (float64, error)
	GetJournalBook(companyID uuid.UUID, start, end time.Time) ([]JournalBookRow, error)
}

type ReportRepository struct {
	Db *gorm.DB
}

func NewReportRepository(db *gorm.DB) IReportRepository {
	return &ReportRepository{Db: db}
}

func (r *ReportRepository) GetLedger(companyID uuid.UUID, start, end time.Time) ([]AccountLedgerRow, error) {
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
		  AND je.company_id = ?
		  AND je.date >= ?
		  AND je.date <= ?
		GROUP BY ca.code, ca.name, cg.name, cs.name
		ORDER BY ca.code
	`, companyID, start, end).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetLedgerUpTo(companyID uuid.UUID, end time.Time) ([]AccountLedgerRow, error) {
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
		  AND je.company_id = ?
		  AND je.date <= ?
		GROUP BY ca.code, ca.name, cg.name, cs.name
		ORDER BY ca.code
	`, companyID, end).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetPostedCount(companyID uuid.UUID, start, end time.Time) (int, error) {
	var count int64
	err := r.Db.Raw(`
		SELECT COUNT(*) FROM journal_entries
		WHERE status = 'posted'
		  AND company_id = ?
		  AND date >= ?
		  AND date <= ?
	`, companyID, start, end).Scan(&count).Error
	return int(count), err
}

// GeneralLedgerRow is a single transaction line for the general ledger report.
type GeneralLedgerRow struct {
	AccountCode   string
	AccountName   string
	GroupName     string
	SubgroupName  string
	Date          string
	JournalNumber string
	Description   string
	Debit         float64
	Credit        float64
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

func (r *ReportRepository) GetGeneralLedger(companyID uuid.UUID, start, end time.Time, coaID string) ([]GeneralLedgerRow, error) {
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
		  AND je.company_id = ?
		  AND je.date >= ?
		  AND je.date <= ?`

	args := []interface{}{companyID, start, end}
	if coaID != "" {
		query += " AND ca.id = ?"
		args = append(args, coaID)
	}
	query += " ORDER BY ca.code, je.date, je.journal_number"

	err := r.Db.Raw(query, args...).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetOpeningBalance(companyID uuid.UUID, asOf time.Time, coaID string) (float64, error) {
	type result struct{ Balance float64 }
	var res result

	query := `
		SELECT COALESCE(SUM(jl.debit) - SUM(jl.credit), 0) AS balance
		FROM journal_lines jl
		INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
		INNER JOIN coa ca             ON ca.id = jl.coa_id
		WHERE je.status = 'posted'
		  AND je.company_id = ?
		  AND je.date < ?`

	args := []interface{}{companyID, asOf}
	if coaID != "" {
		query += " AND ca.code = ?"
		args = append(args, coaID)
	}

	err := r.Db.Raw(query, args...).Scan(&res).Error
	return res.Balance, err
}

func (r *ReportRepository) GetJournalBook(companyID uuid.UUID, start, end time.Time) ([]JournalBookRow, error) {
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
		  AND je.company_id = ?
		  AND je.date >= ?
		  AND je.date <= ?
		ORDER BY je.date, je.journal_number, jl.id
	`, companyID, start, end).Scan(&rows).Error
	return rows, err
}
