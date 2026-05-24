package service

import (
	"math"
	"strings"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/google/uuid"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

func groupType(r repository.AccountLedgerRow) entity.COAGroupType {
	return entity.COAGroupType(r.GroupType)
}

func groupTypeFromGL(r repository.GeneralLedgerRow) entity.COAGroupType {
	return entity.COAGroupType(r.GroupType)
}

// subgroupContains masih dipakai untuk memisahkan sub-kategori dalam satu group type
// (misal: memisahkan "Pendapatan Lain-lain" dari "Pendapatan Usaha" dalam group revenue).
func subgroupContains(name, keyword string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(keyword))
}

type IReportService interface {
	TrialBalance(companyID uuid.UUID, start, end time.Time) (response.TrialBalanceResponse, error)
	ProfitLoss(companyID uuid.UUID, start, end time.Time) (response.ProfitLossResponse, error)
	BalanceSheet(companyID uuid.UUID, asOf time.Time) (response.BalanceSheetResponse, error)
	CashFlow(companyID uuid.UUID, start, end time.Time) (response.CashFlowResponse, error)
	EquityStatement(companyID uuid.UUID, start, end time.Time) (response.EquityStatementResponse, error)
	GeneralLedger(companyID uuid.UUID, start, end time.Time, coaID string) (response.GeneralLedgerResponse, error)
	JournalBook(companyID uuid.UUID, start, end time.Time) (response.JournalBookResponse, error)
}

type ReportService struct {
	IReportRepository repository.IReportRepository
}

func NewReportService(repo repository.IReportRepository) IReportService {
	return &ReportService{IReportRepository: repo}
}

// ─── 1. Trial Balance ────────────────────────────────────────────────────────
func (s *ReportService) TrialBalance(companyID uuid.UUID, start, end time.Time) (response.TrialBalanceResponse, error) {
	rows, err := s.IReportRepository.GetLedger(companyID, start, end)
	if err != nil {
		return response.TrialBalanceResponse{}, err
	}

	postedCount, _ := s.IReportRepository.GetPostedCount(companyID, start, end)

	var lines []response.TrialBalanceLine
	var totalDebit, totalCredit float64

	for _, r := range rows {
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
func (s *ReportService) ProfitLoss(companyID uuid.UUID, start, end time.Time) (response.ProfitLossResponse, error) {
	rows, err := s.IReportRepository.GetLedger(companyID, start, end)
	if err != nil {
		return response.ProfitLossResponse{}, err
	}

	operatingRevenue := response.ReportSection{Title: "Pendapatan Operasional"}
	otherRevenue := response.ReportSection{Title: "Pendapatan Lain-lain"}
	cogs := response.ReportSection{Title: "Harga Pokok Penjualan (HPP)"}
	opex := response.ReportSection{Title: "Beban Operasional"}
	adminExp := response.ReportSection{Title: "Beban Administrasi"}
	finExp := response.ReportSection{Title: "Beban Keuangan / Lain-lain"}

	for _, r := range rows {
		gt := groupType(r)
		sg := strings.ToLower(r.SubgroupName)

		switch gt {
		// ── Revenue ───────────────────────────────────────────────────────────
		case entity.COAGroupTypeRevenue:
			item := response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      r.TotalCredit - r.TotalDebit,
			}
			if subgroupContains(sg, "lain") {
				otherRevenue.Items = append(otherRevenue.Items, item)
				otherRevenue.Subtotal += item.Amount
			} else {
				operatingRevenue.Items = append(operatingRevenue.Items, item)
				operatingRevenue.Subtotal += item.Amount
			}

		// ── COGS / HPP ────────────────────────────────────────────────────────
		case entity.COAGroupTypeCOGS:
			item := response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      r.TotalDebit - r.TotalCredit,
			}
			cogs.Items = append(cogs.Items, item)
			cogs.Subtotal += item.Amount

		// ── Expense / Beban ───────────────────────────────────────────────────
		case entity.COAGroupTypeExpense:
			item := response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      r.TotalDebit - r.TotalCredit,
			}
			switch {
			case subgroupContains(sg, "operasional") || subgroupContains(sg, "operating"):
				opex.Items = append(opex.Items, item)
				opex.Subtotal += item.Amount
			case subgroupContains(sg, "administrasi") || subgroupContains(sg, "admin"):
				adminExp.Items = append(adminExp.Items, item)
				adminExp.Subtotal += item.Amount
			default:
				finExp.Items = append(finExp.Items, item)
				finExp.Subtotal += item.Amount
			}
		}
	}

	totalRevenue := operatingRevenue.Subtotal + otherRevenue.Subtotal
	grossProfit := operatingRevenue.Subtotal - cogs.Subtotal
	totalOpex := opex.Subtotal + adminExp.Subtotal
	operatingProfit := grossProfit - totalOpex
	netProfit := operatingProfit + otherRevenue.Subtotal - finExp.Subtotal

	return response.ProfitLossResponse{
		GeneratedAt:       time.Now(),
		StartDate:         start.Format("2006-01-02"),
		EndDate:           end.Format("2006-01-02"),
		OperatingRevenue:  operatingRevenue,
		OtherRevenue:      otherRevenue,
		TotalRevenue:      totalRevenue,
		COGS:              cogs,
		GrossProfit:       grossProfit,
		OperatingExpenses: opex,
		AdminExpenses:     adminExp,
		FinancialExpenses: finExp,
		OperatingProfit:   operatingProfit,
		NetProfit:         netProfit,
	}, nil
}

// ─── 3. Balance Sheet ────────────────────────────────────────────────────────
func (s *ReportService) BalanceSheet(companyID uuid.UUID, asOf time.Time) (response.BalanceSheetResponse, error) {
	rows, err := s.IReportRepository.GetLedgerUpTo(companyID, asOf)
	if err != nil {
		return response.BalanceSheetResponse{}, err
	}

	assetSections := map[string]*response.ReportSection{}
	liabSections := map[string]*response.ReportSection{}
	equitySections := map[string]*response.ReportSection{}

	var totalAssets, totalLiabilities, totalEquity float64

	for _, r := range rows {
		gt := groupType(r)

		switch gt {
		case entity.COAGroupTypeAsset:
			amount := r.TotalDebit - r.TotalCredit
			sec := ensureSection(assetSections, r.SubgroupName)
			sec.Items = append(sec.Items, response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      amount,
			})
			sec.Subtotal += amount
			totalAssets += amount

		case entity.COAGroupTypeLiability:
			amount := r.TotalCredit - r.TotalDebit
			sec := ensureSection(liabSections, r.SubgroupName)
			sec.Items = append(sec.Items, response.ReportLineItem{
				AccountCode: r.AccountCode,
				AccountName: r.AccountName,
				Amount:      amount,
			})
			sec.Subtotal += amount
			totalLiabilities += amount

		case entity.COAGroupTypeEquity:
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
func (s *ReportService) CashFlow(companyID uuid.UUID, start, end time.Time) (response.CashFlowResponse, error) {
	periodRows, err := s.IReportRepository.GetLedger(companyID, start, end)
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	openingEnd := start.Add(-24 * time.Hour)
	openingRows, err := s.IReportRepository.GetLedgerUpTo(companyID, openingEnd)
	if err != nil {
		return response.CashFlowResponse{}, err
	}

	operating := response.ReportSection{Title: "Aktivitas Operasi"}
	investing := response.ReportSection{Title: "Aktivitas Investasi"}
	financing := response.ReportSection{Title: "Aktivitas Pendanaan"}

	for _, r := range periodRows {
		gt := groupType(r)
		sg := strings.ToLower(r.SubgroupName)

		item := response.ReportLineItem{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
		}

		switch gt {
		case entity.COAGroupTypeRevenue:
			item.Amount = r.TotalCredit - r.TotalDebit
			operating.Items = append(operating.Items, item)
			operating.Subtotal += item.Amount

		case entity.COAGroupTypeCOGS, entity.COAGroupTypeExpense:
			item.Amount = -(r.TotalDebit - r.TotalCredit)
			operating.Items = append(operating.Items, item)
			operating.Subtotal += item.Amount

		case entity.COAGroupTypeAsset:
			// Aset Tidak Lancar / Fixed Assets → Aktivitas Investasi
			if subgroupContains(sg, "tidak lancar") || subgroupContains(sg, "tetap") ||
				subgroupContains(sg, "fixed") || subgroupContains(sg, "kendaraan") ||
				subgroupContains(sg, "peralatan") || subgroupContains(sg, "inventaris") {
				item.Amount = -(r.TotalDebit - r.TotalCredit)
				investing.Items = append(investing.Items, item)
				investing.Subtotal += item.Amount
			}
			// Aset Lancar (kas/bank) tidak dimasukkan langsung ke cashflow,
			// dihitung sebagai opening/closing cash di bawah.

		case entity.COAGroupTypeLiability, entity.COAGroupTypeEquity:
			item.Amount = r.TotalCredit - r.TotalDebit
			financing.Items = append(financing.Items, item)
			financing.Subtotal += item.Amount
		}
	}

	// Opening cash: akun Asset Lancar (kas & bank)
	var openingCash float64
	for _, r := range openingRows {
		if groupType(r) == entity.COAGroupTypeAsset {
			name := strings.ToLower(r.AccountName)
			sg := strings.ToLower(r.SubgroupName)
			if strings.Contains(name, "kas") || strings.Contains(name, "cash") ||
				strings.Contains(name, "bank") || strings.Contains(sg, "lancar") ||
				strings.Contains(sg, "current") {
				openingCash += r.TotalDebit - r.TotalCredit
			}
		}
	}

	netChange := operating.Subtotal + investing.Subtotal + financing.Subtotal
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
func (s *ReportService) EquityStatement(companyID uuid.UUID, start, end time.Time) (response.EquityStatementResponse, error) {
	openingEnd := start.Add(-24 * time.Hour)
	openRows, err := s.IReportRepository.GetLedgerUpTo(companyID, openingEnd)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var openingEquity float64
	for _, r := range openRows {
		if groupType(r) == entity.COAGroupTypeEquity {
			openingEquity += r.TotalCredit - r.TotalDebit
		}
	}

	periodRows, err := s.IReportRepository.GetLedger(companyID, start, end)
	if err != nil {
		return response.EquityStatementResponse{}, err
	}

	var movements []response.EquityMovement
	var equityMovementTotal float64
	for _, r := range periodRows {
		if groupType(r) == entity.COAGroupTypeEquity {
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

	plResp, err := s.ProfitLoss(companyID, start, end)
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

func (s *ReportService) profitLossFromRows(rows []repository.AccountLedgerRow) (response.ProfitLossResponse, error) {
	var totalRevenue, totalExpense float64
	for _, r := range rows {
		gt := groupType(r)
		switch gt {
		case entity.COAGroupTypeRevenue:
			totalRevenue += r.TotalCredit - r.TotalDebit
		case entity.COAGroupTypeCOGS, entity.COAGroupTypeExpense:
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

// ─── 6. General Ledger (Buku Besar) ──────────────────────────────────────────
func (s *ReportService) GeneralLedger(companyID uuid.UUID, start, end time.Time, coaID string) (response.GeneralLedgerResponse, error) {
	rows, err := s.IReportRepository.GetGeneralLedger(companyID, start, end, coaID)
	if err != nil {
		return response.GeneralLedgerResponse{}, err
	}

	type accountKey struct{ code, name, group, subgroup string }
	type accBucket struct {
		key       accountKey
		groupType entity.COAGroupType
		lines     []response.GeneralLedgerLine
		td        float64
		tc        float64
	}

	orderMap := []accountKey{}
	buckets := map[accountKey]*accBucket{}

	for _, r := range rows {
		k := accountKey{r.AccountCode, r.AccountName, r.GroupName, r.SubgroupName}
		if _, ok := buckets[k]; !ok {
			buckets[k] = &accBucket{key: k, groupType: groupTypeFromGL(r)}
			orderMap = append(orderMap, k)
		}
		b := buckets[k]
		b.lines = append(b.lines, response.GeneralLedgerLine{
			Date:          r.Date,
			JournalNumber: r.JournalNumber,
			Description:   r.Description,
			Debit:         r.Debit,
			Credit:        r.Credit,
		})
		b.td += r.Debit
		b.tc += r.Credit
	}

	var accounts []response.GeneralLedgerAccount

	for _, k := range orderMap {
		b := buckets[k]

		// Gunakan IsDebitNormal() dari entity — tidak ada string matching
		isDebitNormal := b.groupType.IsDebitNormal()

		openRaw, _ := s.IReportRepository.GetOpeningBalance(companyID, start, k.code)
		var opening float64
		if isDebitNormal {
			opening = openRaw
		} else {
			opening = -openRaw
		}

		runningBalance := opening
		for i := range b.lines {
			ln := &b.lines[i]
			if isDebitNormal {
				runningBalance += ln.Debit - ln.Credit
			} else {
				runningBalance += ln.Credit - ln.Debit
			}
			ln.Balance = runningBalance
		}

		closingBalance := opening
		if isDebitNormal {
			closingBalance += b.td - b.tc
		} else {
			closingBalance += b.tc - b.td
		}

		accounts = append(accounts, response.GeneralLedgerAccount{
			AccountCode:    k.code,
			AccountName:    k.name,
			GroupName:      k.group,
			SubgroupName:   k.subgroup,
			OpeningBalance: opening,
			Lines:          b.lines,
			TotalDebit:     b.td,
			TotalCredit:    b.tc,
			ClosingBalance: closingBalance,
		})
	}

	return response.GeneralLedgerResponse{
		GeneratedAt: time.Now(),
		StartDate:   start.Format("2006-01-02"),
		EndDate:     end.Format("2006-01-02"),
		Accounts:    accounts,
	}, nil
}

// ─── 7. Journal Book (Jurnal Umum) ───────────────────────────────────────────
func (s *ReportService) JournalBook(companyID uuid.UUID, start, end time.Time) (response.JournalBookResponse, error) {
	rows, err := s.IReportRepository.GetJournalBook(companyID, start, end)
	if err != nil {
		return response.JournalBookResponse{}, err
	}

	type entryKey struct{ date, number, jtype, desc string }
	orderSlice := []entryKey{}
	entryMap := map[entryKey]*response.JournalBookEntry{}

	for _, r := range rows {
		k := entryKey{r.Date, r.JournalNumber, r.JournalType, r.Description}
		if _, ok := entryMap[k]; !ok {
			entryMap[k] = &response.JournalBookEntry{
				Date:          r.Date,
				JournalNumber: r.JournalNumber,
				Type:          r.JournalType,
				Description:   r.Description,
				TotalDebit:    r.TotalDebit,
				TotalCredit:   r.TotalCredit,
			}
			orderSlice = append(orderSlice, k)
		}
		entryMap[k].Lines = append(entryMap[k].Lines, response.JournalBookLine{
			AccountCode: r.AccountCode,
			AccountName: r.AccountName,
			Debit:       r.Debit,
			Credit:      r.Credit,
		})
	}

	var entries []response.JournalBookEntry
	var grandDebit, grandCredit float64
	for _, k := range orderSlice {
		e := entryMap[k]
		entries = append(entries, *e)
		grandDebit += e.TotalDebit
		grandCredit += e.TotalCredit
	}

	return response.JournalBookResponse{
		GeneratedAt: time.Now(),
		StartDate:   start.Format("2006-01-02"),
		EndDate:     end.Format("2006-01-02"),
		Entries:     entries,
		TotalDebit:  grandDebit,
		TotalCredit: grandCredit,
		EntryCount:  len(entries),
	}, nil
}
