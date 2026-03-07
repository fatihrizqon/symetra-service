-- ============================================================
-- CHART OF ACCOUNTS - SEED DATA
-- Groups → Subgroups → COA (fully synced & relational)
-- ============================================================

-- ============================================================
-- 1. COA GROUPS
-- ============================================================
INSERT INTO public.coa_groups (id, company_id, code, name, normal_balance, status)
VALUES
  ('8a7e7b3d-01f5-4d74-9f6c-8f7c2a4b6a01', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '1', 'Assets',      'debit', 1),
  ('b9f9a6c1-02d6-4c5e-8c73-1d6a6a9e0b02', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '2', 'Liabilities',  'credit', 1),
  ('c7e5a4d2-3e0d-4e9b-bd4a-0e5c3a1f7c03', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '3', 'Equity',       'credit', 1),
  ('203d121b-3c8d-48f6-9b41-5032cbe6ffb1', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '4', 'Revenue',      'credit', 1),
  ('e6e16f7e-6534-44f0-bab8-9849465c7046', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '5', 'Expense',      'debit', 1);


-- ============================================================
-- 2. COA SUBGROUPS
--    group_id → FK ke coa_groups.id
-- ============================================================
INSERT INTO public.coa_subgroups (id, company_id, group_id, code, name, status)
VALUES

  -- ASSETS (group: 8a7e7b3d-...)
  ('19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '8a7e7b3d-01f5-4d74-9f6c-8f7c2a4b6a01', '101', 'Current Assets',    1),
  ('0bbc091d-2c4f-4144-91bf-9f8f7dd9a3da', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '8a7e7b3d-01f5-4d74-9f6c-8f7c2a4b6a01', '102', 'Fixed Assets',      1),
  ('687fc218-56d1-45c1-a767-11026dfd5fec', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '8a7e7b3d-01f5-4d74-9f6c-8f7c2a4b6a01', '103', 'Intangible Assets', 1),
  ('dd21ea93-3ff7-4584-b317-feb6155c0c47', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '8a7e7b3d-01f5-4d74-9f6c-8f7c2a4b6a01', '104', 'Other Assets',      1),

  -- LIABILITIES (group: b9f9a6c1-...)
  ('0e7977ba-c52b-441e-b0cd-e4b5396e2cc0', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'b9f9a6c1-02d6-4c5e-8c73-1d6a6a9e0b02', '201', 'Current Liabilities',   1),
  ('ad0d889d-8191-4f43-bb28-a4b2ac33eca0', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'b9f9a6c1-02d6-4c5e-8c73-1d6a6a9e0b02', '202', 'Long Term Liabilities', 1),
  ('bb251bea-7ee9-49d3-880c-8db3af0d4614', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'b9f9a6c1-02d6-4c5e-8c73-1d6a6a9e0b02', '203', 'Other Liabilities',     1),

  -- EQUITY (group: c7e5a4d2-...)
  ('33db7450-6096-4ee1-a1f1-c2dae1aa4704', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c7e5a4d2-3e0d-4e9b-bd4a-0e5c3a1f7c03', '301', 'Owner Equity',      1),
  ('5a03ce54-6f70-4351-8478-89b3c8c46d87', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c7e5a4d2-3e0d-4e9b-bd4a-0e5c3a1f7c03', '302', 'Retained Earnings', 1),
  ('535de875-c302-4c01-ba0a-415cd07b08c0', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c7e5a4d2-3e0d-4e9b-bd4a-0e5c3a1f7c03', '303', 'Drawings',          1),

  -- REVENUE (group: 203d121b-...)
  ('5fc98940-b79d-481a-b27a-047d82a273f8', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '203d121b-3c8d-48f6-9b41-5032cbe6ffb1', '401', 'Operating Revenue', 1),
  ('76fe0d6d-4395-4ef9-917b-5b05262e1c78', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '203d121b-3c8d-48f6-9b41-5032cbe6ffb1', '402', 'Other Revenue',     1),

  -- EXPENSE (group: e6e16f7e-...)
  ('536bb812-78af-4180-b13b-de7f980c6624', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'e6e16f7e-6534-44f0-bab8-9849465c7046', '501', 'Cost of Goods Sold',     1),
  ('730e4e2b-1a4c-4952-993e-0b1dec4cbd4b', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'e6e16f7e-6534-44f0-bab8-9849465c7046', '502', 'Operating Expense',      1),
  ('545f24f5-c47d-400e-88a6-5c64659f9c48', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'e6e16f7e-6534-44f0-bab8-9849465c7046', '503', 'Administrative Expense', 1),
  ('dd89be9d-70f2-4d37-b39e-d97ecc177010', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'e6e16f7e-6534-44f0-bab8-9849465c7046', '504', 'Financial Expense',      1);


-- ============================================================
-- 3. COA (CHART OF ACCOUNTS)
--    subgroup_id → FK ke coa_subgroups.id
-- ============================================================
INSERT INTO public.coa (id, company_id, subgroup_id, code, name, currency_code, active, status)
VALUES

  -- CURRENT ASSETS (101)
  ('c9be7f73-d70c-4f81-be8f-9c721f653987', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '10101', 'Cash',                'IDR', true, 1),
  ('f241874f-6bd5-4c5f-9e79-017297b9dfe0', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '10102', 'Bank',                'IDR', true, 1),
  ('9728ea3d-9692-4f3c-946b-0322a0c4df7b', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '10103', 'Accounts Receivable',  'IDR', true, 1),
  ('fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '10104', 'Inventory',           'IDR', true, 1),
  ('9b78803f-69b0-4d88-924c-f0f8f58593eb', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '19fc40d0-cf06-47b0-bff0-2f9da3cdc428', '10105', 'Prepaid Expenses',    'IDR', true, 1),

  -- FIXED ASSETS (102)
  ('12c94c1d-ed63-4387-948b-590c6d703719', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0bbc091d-2c4f-4144-91bf-9f8f7dd9a3da', '10201', 'Equipment',                'IDR', true, 1),
  ('578541da-339e-4d84-9213-78dd7b63fcef', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0bbc091d-2c4f-4144-91bf-9f8f7dd9a3da', '10202', 'Vehicles',                 'IDR', true, 1),
  ('35b81147-8ada-49c7-8f23-bc142d56401a', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0bbc091d-2c4f-4144-91bf-9f8f7dd9a3da', '10203', 'Buildings',                'IDR', true, 1),
  ('f254867c-210f-4569-8a1c-87201e06a4e3', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0bbc091d-2c4f-4144-91bf-9f8f7dd9a3da', '10204', 'Accumulated Depreciation',  'IDR', true, 1),

  -- INTANGIBLE ASSETS (103)
  ('6106c007-0586-4f95-9ea8-2b9386b9a467', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '687fc218-56d1-45c1-a767-11026dfd5fec', '10301', 'Software', 'IDR', true, 1),
  ('11721263-6ee7-44ec-a0fa-a47add8ec6bc', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '687fc218-56d1-45c1-a767-11026dfd5fec', '10302', 'Licenses', 'IDR', true, 1),

  -- CURRENT LIABILITIES (201)
  ('201f8fb6-8a12-4e16-907f-746a1a3d8b0a', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0e7977ba-c52b-441e-b0cd-e4b5396e2cc0', '20101', 'Accounts Payable', 'IDR', true, 1),
  ('3b856c8e-2d25-4742-8b8d-528540ac4ebf', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0e7977ba-c52b-441e-b0cd-e4b5396e2cc0', '20102', 'Accrued Expenses', 'IDR', true, 1),
  ('46ca12cf-21b5-4f8a-8c1a-93ffc2b28860', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '0e7977ba-c52b-441e-b0cd-e4b5396e2cc0', '20103', 'Tax Payable',      'IDR', true, 1),

  -- LONG TERM LIABILITIES (202)
  ('4ab4c5f2-d2a6-4a6c-a251-d3bbe40efac8', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'ad0d889d-8191-4f43-bb28-a4b2ac33eca0', '20201', 'Bank Loan',     'IDR', true, 1),
  ('7fb99fcf-704d-40d2-a533-de8603591f68', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'ad0d889d-8191-4f43-bb28-a4b2ac33eca0', '20202', 'Bonds Payable', 'IDR', true, 1),

  -- EQUITY (301, 302, 303)
  ('30434024-e2ab-4986-8f92-03c4f22feb3e', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '33db7450-6096-4ee1-a1f1-c2dae1aa4704', '30101', 'Owner Capital',    'IDR', true, 1),
  ('b8b1a1af-becf-43be-be74-21c824f19a98', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '5a03ce54-6f70-4351-8478-89b3c8c46d87', '30201', 'Retained Earnings','IDR', true, 1),
  ('991ec2bc-5b08-4490-94d4-553ecaccda49', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '535de875-c302-4c01-ba0a-415cd07b08c0', '30301', 'Owner Drawings',   'IDR', true, 1),

  -- OPERATING REVENUE (401)
  ('a11275ac-81df-4e8f-82aa-84476e963955', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '5fc98940-b79d-481a-b27a-047d82a273f8', '40101', 'Sales Revenue',   'IDR', true, 1),
  ('4c548201-4c95-4fca-a056-63808b2813e5', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '5fc98940-b79d-481a-b27a-047d82a273f8', '40102', 'Service Revenue', 'IDR', true, 1),

  -- OTHER REVENUE (402)
  ('0a948478-4b70-4597-8ee7-bb94226c058e', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '76fe0d6d-4395-4ef9-917b-5b05262e1c78', '40201', 'Interest Income', 'IDR', true, 1),

  -- COST OF GOODS SOLD (501)
  ('0382ee9c-05fd-47c2-ac52-7002dd374f4f', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '536bb812-78af-4180-b13b-de7f980c6624', '50101', 'Cost of Goods Sold', 'IDR', true, 1),

  -- OPERATING EXPENSE (502)
  ('ff64e581-e5a2-406c-9647-4f7d08c1d447', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '730e4e2b-1a4c-4952-993e-0b1dec4cbd4b', '50201', 'Salary Expense',    'IDR', true, 1),
  ('285a8fba-2941-4848-8cb6-ce2e245609b2', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '730e4e2b-1a4c-4952-993e-0b1dec4cbd4b', '50202', 'Rent Expense',      'IDR', true, 1),
  ('721086f5-8903-4cd6-bbb0-85018755643c', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '730e4e2b-1a4c-4952-993e-0b1dec4cbd4b', '50203', 'Utilities Expense', 'IDR', true, 1),

  -- ADMINISTRATIVE EXPENSE (503)
  ('ed45e579-970f-4f40-b523-5ecf4cea7732', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '545f24f5-c47d-400e-88a6-5c64659f9c48', '50301', 'Office Supplies Expense', 'IDR', true, 1),

  -- FINANCIAL EXPENSE (504)
  ('eb4ae642-6df4-42bc-9e51-e4b0ba5d2df1', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'dd89be9d-70f2-4d37-b39e-d97ecc177010', '50401', 'Bank Charges', 'IDR', true, 1);