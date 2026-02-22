package service

import (
	"math"
	"strings"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/repository"
)

// ─── COA Group name constants (must match your coa_groups.name data) ─────────
// Adjust these if your group names differ.
const (
	groupAssets      = "assets"
	groupLiabilities = "liabilities"
	groupEquity      = "equity"
	groupRevenue     = "revenue"
	groupExpenses    = "expenses"
)

func normalizeGroup(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func isGroup(groupName, target string) bool {
	return strings.Contains(normalizeGroup(groupName), target)
}

// subgroupContains checks if subgroup name contains keyword (case-insensitive)
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

func (s *ReportService) TrialBalance(start, end time.Time) (response.TrialBalanceResponse, error) {
	rows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.TrialBalanceResponse{}, err
	}

	postedCount, _ := s.IReportRepository.GetPostedCount(start, end)

	var lines []response.TrialBalanceLine
	var totalDebit, totalCredit float64

	for _, r := range rows {
		balance := r.TotalDebit - r.TotalCredit
		lines = append(lines, response.TrialBalanceLine{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
			TotalDebit:  r.TotalDebit,
			TotalCredit: r.TotalCredit,
			Balance:     balance,
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

func (s *ReportService) ProfitLoss(start, end time.Time) (response.ProfitLossResponse, error) {
	rows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.ProfitLossResponse{}, err
	}

	revenue := response.ReportSection{Title: "Revenue"}
	cogs := response.ReportSection{Title: "Cost of Goods Sold (HPP)"}
	opex := response.ReportSection{Title: "Operating Expenses"}
	otherExp := response.ReportSection{Title: "Other Expenses"}
	otherRev := response.ReportSection{Title: "Other Revenue"}

	for _, r := range rows {
		g := normalizeGroup(r.GroupName)
		amount := r.TotalCredit - r.TotalDebit // revenue normal balance = credit
		item := response.ReportLineItem{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
		}

		switch {
		case strings.Contains(g, groupRevenue):
			if subgroupContains(r.SubgroupName, "other") {
				item.Amount = r.TotalCredit - r.TotalDebit
				otherRev.Items = append(otherRev.Items, item)
				otherRev.Subtotal += item.Amount
			} else {
				item.Amount = r.TotalCredit - r.TotalDebit
				revenue.Items = append(revenue.Items, item)
				revenue.Subtotal += item.Amount
			}
		case strings.Contains(g, groupExpenses):
			expAmount := r.TotalDebit - r.TotalCredit // expense normal balance = debit
			item.Amount = expAmount
			if subgroupContains(r.SubgroupName, "cogs") || subgroupContains(r.SubgroupName, "cost of goods") || subgroupContains(r.SubgroupName, "hpp") {
				cogs.Items = append(cogs.Items, item)
				cogs.Subtotal += expAmount
			} else if subgroupContains(r.SubgroupName, "other") {
				otherExp.Items = append(otherExp.Items, item)
				otherExp.Subtotal += expAmount
			} else {
				opex.Items = append(opex.Items, item)
				opex.Subtotal += expAmount
			}
		}
		_ = amount
	}

	grossProfit := revenue.Subtotal - cogs.Subtotal
	operatingProfit := grossProfit - opex.Subtotal
	netProfit := operatingProfit + otherRev.Subtotal - otherExp.Subtotal

	return response.ProfitLossResponse{
		GeneratedAt:       time.Now(),
		StartDate:         start.Format("2006-01-02"),
		EndDate:           end.Format("2006-01-02"),
		Revenue:           revenue,
		COGS:              cogs,
		GrossProfit:       grossProfit,
		OperatingExpenses: opex,
		OtherExpenses:     otherExp,
		OtherRevenue:      otherRev,
		OperatingProfit:   operatingProfit,
		NetProfit:         netProfit,
	}, nil
}

// ─── 3. Balance Sheet ────────────────────────────────────────────────────────

func (s *ReportService) BalanceSheet(asOf time.Time) (response.BalanceSheetResponse, error) {
	// Balance sheet is cumulative — from beginning of time to asOf date
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	rows, err := s.IReportRepository.GetLedger(start, asOf)
	if err != nil {
		return response.BalanceSheetResponse{}, err
	}

	// Map subgroup → section items
	assetSections := map[string]*response.ReportSection{}
	liabSections := map[string]*response.ReportSection{}
	equitySections := map[string]*response.ReportSection{}

	var totalAssets, totalLiabilities, totalEquity float64

	// Also compute net profit to add to equity
	plResp, err := s.ProfitLoss(start, asOf)
	if err != nil {
		return response.BalanceSheetResponse{}, err
	}

	for _, r := range rows {
		g := normalizeGroup(r.GroupName)
		item := response.ReportLineItem{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
		}

		switch {
		case strings.Contains(g, groupAssets):
			item.Amount = r.TotalDebit - r.TotalCredit
			sec := ensureSection(assetSections, r.SubgroupName)
			sec.Items = append(sec.Items, item)
			sec.Subtotal += item.Amount
			totalAssets += item.Amount

		case strings.Contains(g, groupLiabilities):
			item.Amount = r.TotalCredit - r.TotalDebit
			sec := ensureSection(liabSections, r.SubgroupName)
			sec.Items = append(sec.Items, item)
			sec.Subtotal += item.Amount
			totalLiabilities += item.Amount

		case strings.Contains(g, groupEquity):
			item.Amount = r.TotalCredit - r.TotalDebit
			sec := ensureSection(equitySections, r.SubgroupName)
			sec.Items = append(sec.Items, item)
			sec.Subtotal += item.Amount
			totalEquity += item.Amount
		}
	}

	// Add Retained Earnings (net profit) to Equity
	if plResp.NetProfit != 0 {
		sec := ensureSection(equitySections, "Retained Earnings")
		sec.Items = append(sec.Items, response.ReportLineItem{
			AccountCode: "-",
			AccountName: "Net Profit (Current Period)",
			Amount:      plResp.NetProfit,
		})
		sec.Subtotal += plResp.NetProfit
		totalEquity += plResp.NetProfit
	}

	assets := sectionsToSlice(assetSections)
	liabilities := sectionsToSlice(liabSections)
	equity := sectionsToSlice(equitySections)

	return response.BalanceSheetResponse{
		GeneratedAt:      time.Now(),
		AsOfDate:         asOf.Format("2006-01-02"),
		Assets:           assets,
		TotalAssets:      totalAssets,
		Liabilities:      liabilities,
		TotalLiabilities: totalLiabilities,
		Equity:           equity,
		TotalEquity:      totalEquity,
		IsBalanced:       math.Abs(totalAssets-(totalLiabilities+totalEquity)) < 0.01,
	}, nil
}

// ─── 4. Cash Flow ────────────────────────────────────────────────────────────

func (s *ReportService) CashFlow(start, end time.Time) (response.CashFlowResponse, error) {
	rows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	// Opening cash: all cash/bank before start date
	openingRows, err := s.IReportRepository.GetLedger(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), start.Add(-24*time.Hour))
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	operating := response.ReportSection{Title: "Operating Activities"}
	investing := response.ReportSection{Title: "Investing Activities"}
	financing := response.ReportSection{Title: "Financing Activities"}

	for _, r := range rows {
		g := normalizeGroup(r.GroupName)
		sg := strings.ToLower(r.SubgroupName)
		netFlow := r.TotalDebit - r.TotalCredit

		item := response.ReportLineItem{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
			Amount:      netFlow,
		}

		switch {
		case strings.Contains(g, groupRevenue) || (strings.Contains(g, groupExpenses) && (strings.Contains(sg, "cogs") || strings.Contains(sg, "operating"))):
			operating.Items = append(operating.Items, item)
			operating.Subtotal += netFlow

		case strings.Contains(g, groupAssets) && (strings.Contains(sg, "fixed") || strings.Contains(sg, "equipment") || strings.Contains(sg, "vehicle")):
			investing.Items = append(investing.Items, item)
			investing.Subtotal += netFlow

		case strings.Contains(g, groupLiabilities) || strings.Contains(g, groupEquity):
			financing.Items = append(financing.Items, item)
			financing.Subtotal += netFlow
		}
	}

	// Opening cash
	var openingCash float64
	for _, r := range openingRows {
		n := strings.ToLower(r.AccountName)
		if strings.Contains(n, "cash") || strings.Contains(n, "bank") {
			openingCash += r.TotalDebit - r.TotalCredit
		}
	}

	netChange := operating.Subtotal + investing.Subtotal + financing.Subtotal
	// Negate because cash normal balance = debit (increase = debit > credit)
	netOperating := -operating.Subtotal
	operating.Subtotal = netOperating

	return response.CashFlowResponse{
		GeneratedAt:      time.Now(),
		StartDate:        start.Format("2006-01-02"),
		EndDate:          end.Format("2006-01-02"),
		Operating:        operating,
		Investing:        investing,
		Financing:        financing,
		NetCashOperating: netOperating,
		NetCashInvesting: investing.Subtotal,
		NetCashFinancing: financing.Subtotal,
		NetCashChange:    netChange,
		OpeningCash:      openingCash,
		ClosingCash:      openingCash + netChange,
	}, nil
}

// ─── 5. Equity Statement ─────────────────────────────────────────────────────

func (s *ReportService) EquityStatement(start, end time.Time) (response.EquityStatementResponse, error) {
	// Opening equity = all equity postings before start
	openStart := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	openRows, err := s.IReportRepository.GetLedger(openStart, start.Add(-24*time.Hour))
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var openingEquity float64
	for _, r := range openRows {
		if strings.Contains(normalizeGroup(r.GroupName), groupEquity) {
			openingEquity += r.TotalCredit - r.TotalDebit
		}
	}

	// Period equity movements
	periodRows, err := s.IReportRepository.GetLedger(start, end)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var movements []response.EquityMovement
	for _, r := range periodRows {
		if strings.Contains(normalizeGroup(r.GroupName), groupEquity) {
			amount := r.TotalCredit - r.TotalDebit
			if amount != 0 {
				movements = append(movements, response.EquityMovement{
					Description: r.AccountName,
					Amount:      amount,
				})
			}
		}
	}

	// Net profit for period
	plResp, err := s.ProfitLoss(start, end)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var equityMovementTotal float64
	for _, m := range movements {
		equityMovementTotal += m.Amount
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

// ─── Helpers ─────────────────────────────────────────────────────────────────

func ensureSection(m map[string]*response.ReportSection, title string) *response.ReportSection {
	if _, ok := m[title]; !ok {
		m[title] = &response.ReportSection{Title: title}
	}
	return m[title]
}

func sectionsToSlice(m map[string]*response.ReportSection) []response.ReportSection {
	var result []response.ReportSection
	for _, sec := range m {
		result = append(result, *sec)
	}
	return result
}
