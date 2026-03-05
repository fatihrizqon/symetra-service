package service

import (
	"math"
	"strings"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/repository"
)

// ─── COA Group name constants ─────────────────────────────────────────────────
// Must match coa_groups.name in seed data (case-insensitive contains check).
const (
	groupAssets      = "assets"
	groupLiabilities = "liabilities"
	groupEquity      = "equity"
	groupRevenue     = "revenue"
	groupExpense     = "expense" // COA seed uses singular "Expense"
)

func normalizeGroup(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func isGroup(groupName, target string) bool {
	return strings.Contains(normalizeGroup(groupName), target)
}

func subgroupContains(name, keyword string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(keyword))
}

type IReportService interface {
	TrialBalance(start, end time.Time) (response.TrialBalanceResponse, error)
	ProfitLoss(start, end time.Time) (response.ProfitLossResponse, error)
	BalanceSheet(asOf time.Time) (response.BalanceSheetResponse, error)
	CashFlow(start, end time.Time) (response.CashFlowResponse, error)
	EquityStatement(start, end time.Time) (response.EquityStatementResponse, error)
}

type ReportService struct {
	IReportRepository repository.IReportRepository
}

func NewReportService(repo repository.IReportRepository) IReportService {
	return &ReportService{IReportRepository: repo}
}

// ─── 1. Trial Balance ────────────────────────────────────────────────────────
// Shows debit / credit totals per account for the given period.

func (s *ReportService) TrialBalance(start, end time.Time) (response.TrialBalanceResponse, error) {
	rows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.TrialBalanceResponse{}, err
	}

	postedCount, _ := s.IReportRepository.GetPostedCount(start, end)

	var lines []response.TrialBalanceLine
	var totalDebit, totalCredit float64

	for _, r := range rows {
		// Skip accounts with zero activity in the period
		if r.TotalDebit == 0 && r.TotalCredit == 0 {
			continue
		}
		balance := r.TotalDebit - r.TotalCredit
		lines = append(lines, response.TrialBalanceLine{
			AccountCode:  r.AccountCode,
			AccountName:  r.AccountName,
			GroupName:    r.GroupName,
			SubgroupName: r.SubgroupName,
			TotalDebit:   r.TotalDebit,
			TotalCredit:  r.TotalCredit,
			Balance:      balance,
		})
		totalDebit += r.TotalDebit
		totalCredit += r.TotalCredit
	}

	return response.TrialBalanceResponse{
		GeneratedAt: time.Now(),
		StartDate:   start.Format("2006-01-02"),
		EndDate:     end.Format("2006-01-02"),
		Lines:       lines,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
		IsBalanced:  math.Abs(totalDebit-totalCredit) < 0.01,
		PostedCount: postedCount,
	}, nil
}

// ─── 2. Profit & Loss ────────────────────────────────────────────────────────
// Laba Rugi: Pendapatan → HPP → Laba Kotor → Beban Operasional →
//            Laba Operasional → Pendapatan/Beban Lain → Laba Bersih

func (s *ReportService) ProfitLoss(start, end time.Time) (response.ProfitLossResponse, error) {
	rows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.ProfitLossResponse{}, err
	}

	operatingRevenue := response.ReportSection{Title: "Pendapatan Operasional"}
	otherRevenue     := response.ReportSection{Title: "Pendapatan Lain-lain"}
	cogs             := response.ReportSection{Title: "Harga Pokok Penjualan (HPP)"}
	opex             := response.ReportSection{Title: "Beban Operasional"}
	adminExp         := response.ReportSection{Title: "Beban Administrasi"}
	finExp           := response.ReportSection{Title: "Beban Keuangan / Lain-lain"}

	for _, r := range rows {
		g  := normalizeGroup(r.GroupName)
		sg := strings.ToLower(r.SubgroupName)

		switch {
		// ── Revenue group ─────────────────────────────────────────────────
		case isGroup(g, groupRevenue):
			item := response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				// Revenue normal balance = credit; positive = credit > debit
				Amount: r.TotalCredit - r.TotalDebit,
			}
			if subgroupContains(sg, "other") {
				otherRevenue.Items = append(otherRevenue.Items, item)
				otherRevenue.Subtotal += item.Amount
			} else {
				operatingRevenue.Items = append(operatingRevenue.Items, item)
				operatingRevenue.Subtotal += item.Amount
			}

		// ── Expense group ─────────────────────────────────────────────────
		case isGroup(g, groupExpense):
			item := response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				// Expense normal balance = debit; positive = debit > credit
				Amount: r.TotalDebit - r.TotalCredit,
			}
			switch {
			case subgroupContains(sg, "cost of goods") || subgroupContains(sg, "cogs") || subgroupContains(sg, "hpp"):
				cogs.Items = append(cogs.Items, item)
				cogs.Subtotal += item.Amount
			case subgroupContains(sg, "operating"):
				opex.Items = append(opex.Items, item)
				opex.Subtotal += item.Amount
			case subgroupContains(sg, "administrative") || subgroupContains(sg, "admin"):
				adminExp.Items = append(adminExp.Items, item)
				adminExp.Subtotal += item.Amount
			default:
				// Financial expense, other → goes to other/financial
				finExp.Items = append(finExp.Items, item)
				finExp.Subtotal += item.Amount
			}
		}
	}

	totalRevenue     := operatingRevenue.Subtotal + otherRevenue.Subtotal
	grossProfit      := operatingRevenue.Subtotal - cogs.Subtotal
	totalOpex        := opex.Subtotal + adminExp.Subtotal
	operatingProfit  := grossProfit - totalOpex
	netProfit        := operatingProfit + otherRevenue.Subtotal - finExp.Subtotal

	return response.ProfitLossResponse{
		GeneratedAt:      time.Now(),
		StartDate:        start.Format("2006-01-02"),
		EndDate:          end.Format("2006-01-02"),
		OperatingRevenue: operatingRevenue,
		OtherRevenue:     otherRevenue,
		TotalRevenue:     totalRevenue,
		COGS:             cogs,
		GrossProfit:      grossProfit,
		OperatingExpenses: opex,
		AdminExpenses:    adminExp,
		FinancialExpenses: finExp,
		OperatingProfit:  operatingProfit,
		NetProfit:        netProfit,
	}, nil
}

// ─── 3. Balance Sheet ────────────────────────────────────────────────────────
// Neraca: kumulatif dari awal sampai asOf.

func (s *ReportService) BalanceSheet(asOf time.Time) (response.BalanceSheetResponse, error) {
	rows, err := s.IReportRepository.GetLedgerUpTo(asOf)
	if err != nil {
		return response.BalanceSheetResponse{}, err
	}

	assetSections  := map[string]*response.ReportSection{}
	liabSections   := map[string]*response.ReportSection{}
	equitySections := map[string]*response.ReportSection{}

	var totalAssets, totalLiabilities, totalEquity float64

	for _, r := range rows {
		g := normalizeGroup(r.GroupName)

		switch {
		case isGroup(g, groupAssets):
			// Asset normal balance = debit
			amount := r.TotalDebit - r.TotalCredit
			sec := ensureSection(assetSections, r.SubgroupName)
			sec.Items = append(sec.Items, response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      amount,
			})
			sec.Subtotal += amount
			totalAssets += amount

		case isGroup(g, groupLiabilities):
			// Liability normal balance = credit
			amount := r.TotalCredit - r.TotalDebit
			sec := ensureSection(liabSections, r.SubgroupName)
			sec.Items = append(sec.Items, response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      amount,
			})
			sec.Subtotal += amount
			totalLiabilities += amount

		case isGroup(g, groupEquity):
			// Equity normal balance = credit
			amount := r.TotalCredit - r.TotalDebit
			sec := ensureSection(equitySections, r.SubgroupName)
			sec.Items = append(sec.Items, response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      amount,
			})
			sec.Subtotal += amount
			totalEquity += amount
		}
	}

	// Retained earnings: net profit from all revenue/expense accounts up to asOf
	// is embedded in the ledger via journal entries; no need to add separately
	// as long as closing entries exist. If no closing entries (open-book),
	// compute and inject current-period net profit into equity.
	plResp, err := s.profitLossFromRows(rows)
	if err == nil && plResp.NetProfit != 0 {
		sec := ensureSection(equitySections, "Laba Periode Berjalan")
		sec.Items = append(sec.Items, response.ReportLineItem{
			AccountCode: "-",
			AccountName: "Laba / Rugi Bersih",
			Amount:      plResp.NetProfit,
		})
		sec.Subtotal += plResp.NetProfit
		totalEquity += plResp.NetProfit
	}

	return response.BalanceSheetResponse{
		GeneratedAt:      time.Now(),
		AsOfDate:         asOf.Format("2006-01-02"),
		Assets:           sectionsToSlice(assetSections),
		TotalAssets:      totalAssets,
		Liabilities:      sectionsToSlice(liabSections),
		TotalLiabilities: totalLiabilities,
		Equity:           sectionsToSlice(equitySections),
		TotalEquity:      totalEquity,
		IsBalanced:       math.Abs(totalAssets-(totalLiabilities+totalEquity)) < 1.0,
	}, nil
}

// ─── 4. Cash Flow ────────────────────────────────────────────────────────────
// Arus Kas: Indirect method — Operating (from P&L), Investing, Financing.

func (s *ReportService) CashFlow(start, end time.Time) (response.CashFlowResponse, error) {
	// Period rows
	periodRows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	// Opening cash: cumulative balance of cash/bank accounts BEFORE start
	openingEnd := start.Add(-24 * time.Hour)
	openingRows, err := s.IReportRepository.GetLedgerUpTo(openingEnd)
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	operating := response.ReportSection{Title: "Aktivitas Operasi"}
	investing  := response.ReportSection{Title: "Aktivitas Investasi"}
	financing  := response.ReportSection{Title: "Aktivitas Pendanaan"}

	for _, r := range periodRows {
		g  := normalizeGroup(r.GroupName)
		sg := strings.ToLower(r.SubgroupName)

		item := response.ReportLineItem{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
		}

		switch {
		// Revenue → cash inflow (positive)
		case isGroup(g, groupRevenue):
			item.Amount = r.TotalCredit - r.TotalDebit
			operating.Items = append(operating.Items, item)
			operating.Subtotal += item.Amount

		// Expense → cash outflow (negative)
		case isGroup(g, groupExpense):
			item.Amount = -(r.TotalDebit - r.TotalCredit)
			operating.Items = append(operating.Items, item)
			operating.Subtotal += item.Amount

		// Fixed / intangible assets → investing
		case isGroup(g, groupAssets) &&
			(subgroupContains(sg, "fixed") || subgroupContains(sg, "intangible") ||
				subgroupContains(sg, "equipment") || subgroupContains(sg, "vehicle")):
			item.Amount = -(r.TotalDebit - r.TotalCredit) // purchase = outflow
			investing.Items = append(investing.Items, item)
			investing.Subtotal += item.Amount

		// Liabilities & equity → financing
		case isGroup(g, groupLiabilities) || isGroup(g, groupEquity):
			item.Amount = r.TotalCredit - r.TotalDebit
			financing.Items = append(financing.Items, item)
			financing.Subtotal += item.Amount
		}
	}

	// Compute opening cash balance (Cash + Bank accounts)
	var openingCash float64
	for _, r := range openingRows {
		if isGroup(normalizeGroup(r.GroupName), groupAssets) {
			name := strings.ToLower(r.AccountName)
			sg   := strings.ToLower(r.SubgroupName)
			if strings.Contains(name, "cash") || strings.Contains(name, "bank") ||
				strings.Contains(sg, "cash") || strings.Contains(sg, "current") {
				openingCash += r.TotalDebit - r.TotalCredit
			}
		}
	}

	netChange   := operating.Subtotal + investing.Subtotal + financing.Subtotal
	closingCash := openingCash + netChange

	return response.CashFlowResponse{
		GeneratedAt:      time.Now(),
		StartDate:        start.Format("2006-01-02"),
		EndDate:          end.Format("2006-01-02"),
		Operating:        operating,
		Investing:        investing,
		Financing:        financing,
		NetCashOperating: operating.Subtotal,
		NetCashInvesting: investing.Subtotal,
		NetCashFinancing: financing.Subtotal,
		NetCashChange:    netChange,
		OpeningCash:      openingCash,
		ClosingCash:      closingCash,
	}, nil
}

// ─── 5. Equity Statement ─────────────────────────────────────────────────────

func (s *ReportService) EquityStatement(start, end time.Time) (response.EquityStatementResponse, error) {
	// Opening equity: cumulative up to day before start
	openingEnd := start.Add(-24 * time.Hour)
	openRows, err := s.IReportRepository.GetLedgerUpTo(openingEnd)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var openingEquity float64
	for _, r := range openRows {
		if isGroup(normalizeGroup(r.GroupName), groupEquity) {
			openingEquity += r.TotalCredit - r.TotalDebit
		}
	}

	// Period equity movements
	periodRows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var movements []response.EquityMovement
	var equityMovementTotal float64
	for _, r := range periodRows {
		if isGroup(normalizeGroup(r.GroupName), groupEquity) {
			amount := r.TotalCredit - r.TotalDebit
			if amount != 0 {
				movements = append(movements, response.EquityMovement{
					AccountCode: r.AccountCode,
					Description: r.AccountName,
					Amount:      amount,
				})
				equityMovementTotal += amount
			}
		}
	}

	// Net profit for period
	plResp, err := s.ProfitLoss(start, end)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	closingEquity := openingEquity + equityMovementTotal + plResp.NetProfit

	return response.EquityStatementResponse{
		GeneratedAt:   time.Now(),
		StartDate:     start.Format("2006-01-02"),
		EndDate:       end.Format("2006-01-02"),
		OpeningEquity: openingEquity,
		Movements:     movements,
		NetProfit:     plResp.NetProfit,
		ClosingEquity: closingEquity,
	}, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// profitLossFromRows computes net profit directly from a pre-fetched ledger slice.
// Used by BalanceSheet to avoid a second DB round-trip.
func (s *ReportService) profitLossFromRows(rows []repository.AccountLedgerRow) (response.ProfitLossResponse, error) {
	var totalRevenue, totalExpense float64
	for _, r := range rows {
		g := normalizeGroup(r.GroupName)
		switch {
		case isGroup(g, groupRevenue):
			totalRevenue += r.TotalCredit - r.TotalDebit
		case isGroup(g, groupExpense):
			totalExpense += r.TotalDebit - r.TotalCredit
		}
	}
	return response.ProfitLossResponse{NetProfit: totalRevenue - totalExpense}, nil
}

func ensureSection(m map[string]*response.ReportSection, title string) *response.ReportSection {
	if _, ok := m[title]; !ok {
		m[title] = &response.ReportSection{Title: title}
	}
	return m[title]
}

func sectionsToSlice(m map[string]*response.ReportSection) []response.ReportSection {
	result := make([]response.ReportSection, 0, len(m))
	for _, sec := range m {
		result = append(result, *sec)
	}
	return result
}
