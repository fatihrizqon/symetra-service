# Symetra Enterprise Upgrade — Implementation Summary

## Architecture Overview
```
PostgreSQL (GORM) ← Service Layer ← Handler Layer ← Fiber Router
                       ↑
              Posting Engine (double-entry enforcement)
```

## Backend: symetra-service-new

### Entity Layer (`internal/entity/`)
| File | Entities |
|------|---------|
| tenant.go | Tenant |
| company.go | Company |
| branch.go | Branch, CostCenter |
| coa.go | COAGroup (refactored), COASubGroup, COA + company_id |
| fiscal_period.go | FiscalYear, FiscalPeriod (OPEN/CLOSED/LOCKED) |
| ledger.go | Ledger, PostingRule, LedgerEntry (**immutable**), LedgerEntryDetail |
| audit.go | WorkflowLog, AuditLog (append-only) |
| rbac.go | Role, RolePermission, UserRole |
| tax.go | TaxCode, VATOutputRegister, VATInputRegister, WHTRegister |
| master_data.go | Customer, Vendor, Item, Warehouse |
| sales.go | SalesQuotation→Order→DeliveryOrder→Invoice→Payment→Return |
| purchase.go | PurchaseRequest→Order→GoodsReceipt→VendorBill→Payment→Return |
| inventory.go | StockMovement, ValuationLayer, StockOpname, BankAccount, BankStatement |

### Repository Layer (`internal/repository/`)
- tenant, company, fiscal_period (FiscalYear + FiscalPeriod)
- ledger (Ledger, PostingRule, LedgerEntry with idempotency)
- sales_invoice, vendor_bill, customer, vendor, tax_code

### Service Layer (`internal/service/`)
| Service | Key Logic |
|---------|-----------|
| **posting_engine** | Double-entry enforcement, ValidatePeriod, Reverse, idempotency |
| **sales_invoice** | DRAFT→SUBMITTED→APPROVED→POSTED→VOID workflow |
| **vendor_bill** | Same workflow + VAT input register creation on POST |
| **customer_payment** | Validates no overpayment, updates invoice outstanding |
| **fiscal_year** | Create + auto-generate 12 monthly periods |
| **fiscal_period** | Status transitions (OPEN/CLOSED/LOCKED — LOCKED is terminal) |
| **enterprise_report** | TrialBalance, BalanceSheet, P&L, AR/AP Aging, VAT Report |

### API Endpoints (`/api/v1/`)
```
# Fiscal
POST   /fiscal-years
GET    /fiscal-years, /fiscal-years/:id
GET    /fiscal-periods
PATCH  /fiscal-periods/:id/status  { status: "CLOSED"|"LOCKED" }

# Sales
POST   /sales-invoices
GET    /sales-invoices, /sales-invoices/:id
PATCH  /sales-invoices/:id/action  { action: "submit"|"approve"|"reject"|"post"|"void" }

POST   /customer-payments
GET    /customer-payments, /customer-payments/:id
PATCH  /customer-payments/:id/action

# Purchase
POST   /vendor-bills
GET    /vendor-bills, /vendor-bills/:id
PATCH  /vendor-bills/:id/action

# Reports
GET    /reports/trial-balance?company_id=&period_id=&hide_zero=true
GET    /reports/balance-sheet?company_id=&period_id=
GET    /reports/profit-loss?company_id=&period_id=&ytd=true
GET    /reports/ar-aging?company_id=&as_of_date=2025-12-31
GET    /reports/ap-aging?company_id=&as_of_date=2025-12-31
GET    /reports/vat?company_id=&period_id=
```

### Multi-Tenancy
- All transactional entities carry `company_id`
- `X-Company-ID` header required on all authenticated endpoints
- UserRole scoped to company (nil = super admin)

### Accounting Integrity
- **Immutability**: LedgerEntry/LedgerEntryDetail — no UPDATE/DELETE
- **Reversal**: Mirror-image entries (debit↔credit swapped)
- **Idempotency**: Unique(reference_type, reference_id, ledger_id)
- **Period control**: CLOSED rejects new posts (409), LOCKED is irreversible
- **Double-entry**: SUM(debit) = SUM(credit) enforced with 0.001 tolerance
- **Balance sheet integrity**: Assets = Liabilities + Equity alert on mismatch

## Frontend: symetra-ui-new

### New Modules (`src/modules/dashboard/`)
| Route | Module | Description |
|-------|--------|-------------|
| `/dashboard/fiscal-periods` | fiscal-years | FiscalYear list + period status management |
| `/dashboard/items` | items | Item/service master data |
| `/dashboard/tax-codes` | tax-codes | VAT/WHT tax code management |
| `/dashboard/sales-invoices` | sales-invoices | Invoice list + detail + workflow actions |
| `/dashboard/customer-payments` | customer-payments | Payment list + detail + workflow |
| `/dashboard/vendor-bills` | vendor-bills | Vendor bill list + detail + workflow |
| `/dashboard/financial-reports` | financial-reports | Balance Sheet, P&L, Trial Balance, AR/AP Aging |
| `/dashboard/tax-reports` | tax-reports | VAT report (output/input register) |

### Shared Types (`src/types/enterprise.ts`)
- Company, FiscalYear, FiscalPeriod, TaxCode, Item
- SalesInvoice + SalesInvoiceLine
- VendorBill + VendorBillLine
- CustomerPayment
- STATUS_COLORS map, WORKFLOW_ACTIONS map

## Database Seed (`database/seed.go`)
Auto-seeds on startup (idempotent):
1. Tenant: "Demo Organization"
2. Company: "PT Demo Utama" (IDR, NPWP: 01.234.567.8-901.000)
3. FiscalYear + 12 monthly periods (OPEN)
4. 4 Ledgers (ACCRUAL, CASH, TAX, MANAGEMENT)
5. COA: 6 groups, 20+ accounts (Kas, Bank, Piutang, Hutang, Penjualan, HPP, Beban)
6. Tax Codes: PPN 11%, PPh23 2% (Jasa), PPh23 15% (Dividen)
7. Admin user: `admin@symetra.id` / `Admin@123`

## Pending (Next Sprint)
- [ ] Posting rules configuration per transaction type
- [ ] Stock movement service + inventory valuation
- [ ] Bank reconciliation service
- [ ] Purchase order → goods receipt → vendor bill 3-way matching
- [ ] Goods delivery order → delivery note → sales invoice auto-create
- [ ] RBAC middleware enforcement on all endpoints
- [ ] Jest unit tests (280 test cases planned)
- [ ] k6 load tests
- [ ] React forms for creating sales invoices, vendor bills, payments
- [ ] Coretax export integration
