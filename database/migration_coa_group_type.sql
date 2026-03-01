-- ============================================================================
-- Migration: Add `type` column to coa_groups
-- Run this ONCE on your existing database before deploying the new backend.
-- ============================================================================

-- Step 1: Add the column with a safe default
ALTER TABLE coa_groups
  ADD COLUMN IF NOT EXISTS type VARCHAR NOT NULL DEFAULT 'asset';

-- Step 2: Set sensible defaults based on normal_balance.
-- Accounts with credit normal balance are most likely revenue or liability.
-- Accounts with debit normal balance are most likely asset or expense.
-- You MUST review and correct these after running the script.
UPDATE coa_groups
SET type = CASE
  WHEN LOWER(normal_balance) = 'credit' THEN 'liability'
  WHEN LOWER(normal_balance) = 'debit'  THEN 'asset'
  ELSE 'asset'
END;

-- Step 3: MANUALLY correct each group to its proper type.
-- Run these UPDATE statements based on your actual COA group data.
-- Example (adjust the names/codes to match your actual data):
--
-- Revenue groups (pendapatan):
-- UPDATE coa_groups SET type = 'revenue' WHERE name ILIKE '%pendapatan%';
-- UPDATE coa_groups SET type = 'revenue' WHERE code IN ('4', '40', '41');
--
-- Expense groups (beban / biaya):
-- UPDATE coa_groups SET type = 'expense' WHERE name ILIKE '%beban%';
-- UPDATE coa_groups SET type = 'expense' WHERE name ILIKE '%biaya%';
-- UPDATE coa_groups SET type = 'expense' WHERE code IN ('5', '50', '51');
--
-- Equity groups:
-- UPDATE coa_groups SET type = 'equity' WHERE name ILIKE '%modal%';
-- UPDATE coa_groups SET type = 'equity' WHERE name ILIKE '%ekuitas%';
--
-- Asset groups (already defaulted above, but confirm):
-- UPDATE coa_groups SET type = 'asset' WHERE name ILIKE '%aset%';
-- UPDATE coa_groups SET type = 'asset' WHERE name ILIKE '%kas%';
--
-- Liability groups:
-- UPDATE coa_groups SET type = 'liability' WHERE name ILIKE '%hutang%';
-- UPDATE coa_groups SET type = 'liability' WHERE name ILIKE '%liabilit%';

-- Step 4: Verify the result
SELECT id, code, name, normal_balance, type FROM coa_groups ORDER BY code;

-- Step 5: Once you're happy with the data, the Go AutoMigrate will handle
-- the column going forward. No further SQL needed.
