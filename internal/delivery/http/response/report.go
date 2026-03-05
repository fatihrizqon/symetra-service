package response

import "time"

// ─── Shared ──────────────────────────────────────────────────────────────────

type ReportLineItem struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

type ReportSection struct {
	Title    string           `json:"title"`
	Items    []ReportLineItem `json:"items"`
	Subtotal float64          `json:"subtotal"`
}

// ─── 1. Trial Balance ────────────────────────────────────────────────────────

type TrialBalanceLine struct {
	AccountCode  string  `json:"account_code"`
	AccountName  string  `json:"account_name"`
	GroupName    string  `json:"group_name"`
	SubgroupName string  `json:"subgroup_name"`
	TotalDebit   float64 `json:"total_debit"`
	TotalCredit  float64 `json:"total_credit"`
	Balance      float64 `json:"balance"`
}

type TrialBalanceResponse struct {
	GeneratedAt time.Time          `json:"generated_at"`
	StartDate   string             `json:"start_date"`
	EndDate     string             `json:"end_date"`
	Lines       []TrialBalanceLine `json:"lines"`
	TotalDebit  float64            `json:"total_debit"`
	TotalCredit float64            `json:"total_credit"`
	IsBalanced  bool               `json:"is_balanced"`
	PostedCount int                `json:"posted_count"`
}

// ─── 2. Profit & Loss ────────────────────────────────────────────────────────
// Laba Rugi:
//   Pendapatan Operasional
//   - HPP
//   = Laba Kotor
//   - Beban Operasional
//   - Beban Administrasi
//   = Laba Operasional
//   + Pendapatan Lain-lain
//   - Beban Keuangan / Lain-lain
//   = Laba Bersih

type ProfitLossResponse struct {
	GeneratedAt       time.Time     `json:"generated_at"`
	StartDate         string        `json:"start_date"`
	EndDate           string        `json:"end_date"`
	OperatingRevenue  ReportSection `json:"operating_revenue"`
	OtherRevenue      ReportSection `json:"other_revenue"`
	TotalRevenue      float64       `json:"total_revenue"`
	COGS              ReportSection `json:"cogs"`
	GrossProfit       float64       `json:"gross_profit"`
	OperatingExpenses ReportSection `json:"operating_expenses"`
	AdminExpenses     ReportSection `json:"admin_expenses"`
	FinancialExpenses ReportSection `json:"financial_expenses"`
	OperatingProfit   float64       `json:"operating_profit"`
	NetProfit         float64       `json:"net_profit"`
}

// ─── 3. Balance Sheet ────────────────────────────────────────────────────────

type BalanceSheetResponse struct {
	GeneratedAt      time.Time       `json:"generated_at"`
	AsOfDate         string          `json:"as_of_date"`
	Assets           []ReportSection `json:"assets"`
	TotalAssets      float64         `json:"total_assets"`
	Liabilities      []ReportSection `json:"liabilities"`
	TotalLiabilities float64         `json:"total_liabilities"`
	Equity           []ReportSection `json:"equity"`
	TotalEquity      float64         `json:"total_equity"`
	IsBalanced       bool            `json:"is_balanced"`
}

// ─── 4. Cash Flow ────────────────────────────────────────────────────────────

type CashFlowResponse struct {
	GeneratedAt      time.Time     `json:"generated_at"`
	StartDate        string        `json:"start_date"`
	EndDate          string        `json:"end_date"`
	Operating        ReportSection `json:"operating"`
	Investing        ReportSection `json:"investing"`
	Financing        ReportSection `json:"financing"`
	NetCashOperating float64       `json:"net_cash_operating"`
	NetCashInvesting float64       `json:"net_cash_investing"`
	NetCashFinancing float64       `json:"net_cash_financing"`
	NetCashChange    float64       `json:"net_cash_change"`
	OpeningCash      float64       `json:"opening_cash"`
	ClosingCash      float64       `json:"closing_cash"`
}

// ─── 5. Equity Statement ─────────────────────────────────────────────────────

type EquityMovement struct {
	AccountCode string  `json:"account_code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type EquityStatementResponse struct {
	GeneratedAt   time.Time        `json:"generated_at"`
	StartDate     string           `json:"start_date"`
	EndDate       string           `json:"end_date"`
	OpeningEquity float64          `json:"opening_equity"`
	Movements     []EquityMovement `json:"movements"`
	NetProfit     float64          `json:"net_profit"`
	ClosingEquity float64          `json:"closing_equity"`
}
