-- ============================================================
-- JURNAL JANUARI 2026 — PT. Global Putra Ventura Sejahtera
-- Generated: 2026-04-05 14:15:19
-- Total entries: 10
-- ============================================================

BEGIN;

-- [JNL-JAN26-001] Data Sales - ISAT GPS (paid - inv des 2025)
-- Debit: 164,175,508.26  |  Kredit: 164,175,508.26
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '50a6e03d-12d0-453c-807b-192fb03be404',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-001',
  '2026-01-31',
  'Data Sales - ISAT GPS (paid - inv des 2025)',
  'general',
  'posted',
  164175508.26,
  164175508.26,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '22fea9a7-1a30-41df-a1ec-1bd4708077de',
  '50a6e03d-12d0-453c-807b-192fb03be404',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  161217391.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '8afddbbb-5278-421a-b3c7-1b8654dd588e',
  '50a6e03d-12d0-453c-807b-192fb03be404',
  '944d01a8-127f-4abe-9346-2b56105945be',
  'Piutang Usaha IDR - Indosat Tbk',
  0.00,
  164175507.93,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '1873ce9a-1e0a-4095-9307-b9b5e33542b2',
  '50a6e03d-12d0-453c-807b-192fb03be404',
  '3f481d2b-2874-4353-9ed7-65fac2526c02',
  'PPh 23 Penjualan',
  2958117.26,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '859f3ca5-116c-4f9b-9c7a-afd219401cc1',
  '50a6e03d-12d0-453c-807b-192fb03be404',
  'a9f8cf8f-63b3-419d-80d7-c7b7147761ba',
  'Pendapatan Diluar Usaha Lainnya',
  0.00,
  0.33,
  NOW(),
  NOW()
);

-- [JNL-JAN26-002] Data Sales - Nastiti (paid & unpaid)
-- Debit: 19,400,000.00  |  Kredit: 19,400,000.00
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '6682e612-a9b8-4f61-8745-3a0bdaf8d4ab',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-002',
  '2026-01-31',
  'Data Sales - Nastiti (paid & unpaid)',
  'general',
  'posted',
  19400000.00,
  19400000.00,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'd34b84a0-646c-4936-9f7e-09369082ab3b',
  '6682e612-a9b8-4f61-8745-3a0bdaf8d4ab',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  18050000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '929651ca-96ec-495d-b6e7-225efe5ba2f7',
  '6682e612-a9b8-4f61-8745-3a0bdaf8d4ab',
  '5ad6ca03-4134-4d5c-bb3d-ef63bb8b4c90',
  'Penjualan PT. GPS - Nastiti',
  0.00,
  19400000.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'c83a6f0f-19fe-4999-a0e8-45f6bcb007fa',
  '6682e612-a9b8-4f61-8745-3a0bdaf8d4ab',
  'c6b6fa22-cb69-4737-bc8a-854cdceebce7',
  'Piutang Usaha IDR - Nastiti',
  1350000.00,
  0.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-003] Data Sales - ISAT GPS (jan 2026)
-- Debit: 334,634,308.50  |  Kredit: 334,634,308.50
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '0878e3f3-a8d2-4ba2-a8ab-4ed795831a86',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-003',
  '2026-01-31',
  'Data Sales - ISAT GPS (jan 2026)',
  'general',
  'posted',
  334634308.50,
  334634308.50,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'a73f0a2b-19eb-4ad6-a865-5cdbb53aad00',
  '0878e3f3-a8d2-4ba2-a8ab-4ed795831a86',
  '944d01a8-127f-4abe-9346-2b56105945be',
  'Piutang Usaha IDR - Indosat Tbk',
  334634308.50,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '5df3dd7e-7c29-4a13-afc7-2f3edc95506d',
  '0878e3f3-a8d2-4ba2-a8ab-4ed795831a86',
  '5ad6ca03-4134-4d5c-bb3d-ef63bb8b4c90',
  'Penjualan PT. GPS - Indosat tbk',
  0.00,
  301472350.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '3563246c-68ea-4ec9-aeaf-e5b22918fd3f',
  '0878e3f3-a8d2-4ba2-a8ab-4ed795831a86',
  '728ea14d-f19d-4a85-96f3-28893598a08a',
  'Hutang PPN',
  0.00,
  33161958.50,
  NOW(),
  NOW()
);

-- [JNL-JAN26-004] Data Biaya
-- Debit: 69,836,016.79  |  Kredit: 69,836,016.79
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-004',
  '2026-01-31',
  'Data Biaya',
  'general',
  'posted',
  69836016.79,
  69836016.79,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '80f5b0e9-b54e-474d-a7ac-1aed72da3983',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Adm. Bank & Buku Cek/Giro',
  181093.30,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '3b0fb1e1-b8bd-4c77-9102-17409fbaea4e',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '5f766cdf-7933-4616-9386-26fd674208cf',
  'Beban Bunga Pinjaman',
  1680000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'e7215095-f3d1-4025-b660-8348bb7d70a4',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Entertaint',
  200200.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'f9525244-c446-4805-b0ef-23a25cf9d24f',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  'e92061f9-3e14-4333-83b6-751ae4ecd3a0',
  'Beban Gaji, Tunjangan',
  1500000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '1a082963-1d24-4568-86d3-1769f2391dc5',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  'e92061f9-3e14-4333-83b6-751ae4ecd3a0',
  'Beban Gaji, Tunjangan (Lapangan)',
  500000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'b69fc074-6274-4d07-88d7-f0ea091f16ae',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Konsumsi',
  500000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'db8c3d9a-6d83-450f-a60a-d991ebf29f0b',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Listrik',
  1158735.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '435e3200-ef52-4929-bfdc-d6763a0da195',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Operasional Lainnya',
  2377084.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'b29e6af1-51aa-4e62-95b5-578651f69cd1',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Retribusi & Sumbangan',
  100000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'da5b73c3-56de-499a-ad8f-7d2bb5f81a9d',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Rumah Tangga/Pantry',
  125000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '45190409-06ee-4cb8-9ebc-caffbcc5abf1',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '1a3bcab7-ee01-4cda-88ab-deb5919800a0',
  'Beban Sewa Gedung',
  2500000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'dd5a3427-aad9-479c-a40c-48578be1715e',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Telekomunikasi',
  35000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '207d3aad-c58d-4cd5-be4e-a5dfa357fcba',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Transportasi',
  450000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '9cfaa68a-ae27-4864-bfec-f4a2e105d826',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Tunjangan Kesehatan',
  5000000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '72706ba5-4785-4f59-b02e-93993b7fe7a8',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Biaya Operasional Lapangan',
  15752824.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '86887b35-14b9-46c9-bf6c-f143d8462dfc',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Biaya Operasional Lapangan (FL)',
  37775300.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '410e963c-88ec-4396-bdb1-83f25ac57921',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '8c93fb5c-9bfe-4901-b97c-da3b793a55f6',
  'Pajak Jasa Giro',
  780.49,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '57c011fe-ddde-42ca-bc64-472b309ba939',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  0.00,
  69795923.49,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'fc6a9528-a0f4-491b-9644-37921cd1097d',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  'be82b7c6-5e92-4904-b479-2f196d5de945',
  'Kas & Bank - Bank Mandiri GPVS',
  0.00,
  93.30,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '72d02ede-9cda-4b18-94a8-fd0c423aff5d',
  'a7fe1c1e-a9cd-4617-8868-8e14e45fabe9',
  'b77d370d-696c-421d-8fe9-57f51ed0a99b',
  'Kas & Bank - Bank BSI',
  0.00,
  40000.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-005] Data Bank-Kas
-- Debit: 2,100,000.00  |  Kredit: 2,100,000.00
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '83c0bab8-9808-4af1-a11f-a6f210ecc4a9',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-005',
  '2026-01-31',
  'Data Bank-Kas',
  'general',
  'posted',
  2100000.00,
  2100000.00,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '022685bb-3246-4196-99ac-e0b570510e2a',
  '83c0bab8-9808-4af1-a11f-a6f210ecc4a9',
  '0c59162d-be10-4517-9eeb-ebf9b9237b6d',
  'Kas & Bank - Kas Kecil Rizky',
  2100000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '87e62353-a8ca-4d68-93b0-72dfc8eddb67',
  '83c0bab8-9808-4af1-a11f-a6f210ecc4a9',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  0.00,
  2100000.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-006] Data Transaksi Lainnya
-- Debit: 151,587,381.52  |  Kredit: 151,587,381.52
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-006',
  '2026-01-31',
  'Data Transaksi Lainnya',
  'general',
  'posted',
  151587381.52,
  151587381.52,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '46d64b39-9b4a-4528-9303-3c3b7f550487',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  '5f766cdf-7933-4616-9386-26fd674208cf',
  'Beban Bunga Pinjaman',
  2698351.52,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '309f8b6e-4493-4223-a94f-998e2d4bb6d4',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  '04a1a6e0-0708-4a50-9a4b-8ecfa5cbf462',
  'KPPS Tamzis Bina Utama (1)',
  15083333.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '34c8a75a-dd99-418f-ba59-f7b456b525ca',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  '71ef9c61-0bb3-428f-b709-fb3dd35b1cce',
  'Piutang Karyawan',
  2500000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '699286b6-22e6-4236-9a79-7cfb2ac1942c',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  '2fdaa297-84e0-42bc-8c9b-7c9d9e5973d8',
  'Piutang - Nur Anita Meikhawati',
  5000000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'ce4cf418-b68e-4933-91dd-9ba78ecc0a79',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  'cc482102-fe6b-4674-ba19-896584a1546d',
  'Piutang Direksi - Nanang Ranu',
  6000000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '64bc3089-22c1-46ae-91f8-2dc4acaf80f8',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  'cc54c6f5-c155-4a3c-aa75-cce20b31e9b5',
  'Hutang Gadai (Kendaraan)',
  4347500.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'e06557b7-5ac5-482a-a841-c5b31a678f85',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  'a2476a89-5505-4fb9-9607-3cb58adc66f6',
  'Utang Usaha',
  115958197.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '065fb3d6-1a04-4012-9222-7ace938c61b2',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  'b77d370d-696c-421d-8fe9-57f51ed0a99b',
  'Kas & Bank - Bank BSI',
  0.00,
  17781684.52,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'ad80a18b-b3fc-4f1d-9625-536d76cbad3a',
  '9c6fec69-5bf5-47a5-9c58-789297ec89f4',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  0.00,
  133805697.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-007] Data Gaji (Des)
-- Debit: 81,375,000.00  |  Kredit: 81,375,000.00
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '32e1245f-142b-4436-a78d-bd7c09b81424',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-007',
  '2026-01-31',
  'Data Gaji (Des)',
  'general',
  'posted',
  81375000.00,
  81375000.00,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '21011d24-dfd0-4f22-aa40-37f66eeb623a',
  '32e1245f-142b-4436-a78d-bd7c09b81424',
  'e92061f9-3e14-4333-83b6-751ae4ecd3a0',
  'Beban Gaji, Tunjangan',
  44300000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'd9664b8d-2eed-4283-b070-544e3e9b6bb9',
  '32e1245f-142b-4436-a78d-bd7c09b81424',
  'e92061f9-3e14-4333-83b6-751ae4ecd3a0',
  'Beban Gaji, Tunjangan (Lapangan)',
  37075000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '14e2ff4d-1f1c-4c21-8c2e-3c874d57528c',
  '32e1245f-142b-4436-a78d-bd7c09b81424',
  '9cd442d8-b99e-48a9-9d91-d971f41083e1',
  'Kas & Bank - Bank Mandiri GPS',
  0.00,
  79875000.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '655241e4-fff2-4de9-803c-e2bad2cb9557',
  '32e1245f-142b-4436-a78d-bd7c09b81424',
  '71ef9c61-0bb3-428f-b709-fb3dd35b1cce',
  'Piutang Karyawan',
  0.00,
  1500000.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-008] Data Petty Cash
-- Debit: 40,870,390.00  |  Kredit: 40,870,390.00
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-008',
  '2026-01-31',
  'Data Petty Cash',
  'general',
  'posted',
  40870390.00,
  40870390.00,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '6e1a4670-5128-4182-8741-ee2e4e360134',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Alat Tulis Kantor',
  234000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'd1f27a1d-6d0c-468e-b039-f90f39004dbe',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Internet',
  565420.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '95e392ab-f38a-41ff-9f8d-1847865680a2',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Jasa Profesional',
  1000000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'b191943f-ce7f-4747-9763-a23bf9531dbf',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Konsumsi',
  863000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '19714ed1-63e1-48a4-ae26-a5de912c581e',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban PAM',
  89970.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '656bbc24-e01e-4133-a09a-7f0a6f2dc63f',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Retribusi & Sumbangan',
  150000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'fac51e1f-d77c-4d67-ad7d-83bc43b5442a',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Rumah Tangga/Pantry',
  1875500.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '644cc878-2dd7-47bf-b66f-e4624a2accd0',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Beban Telekomunikasi',
  50000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '755a5e50-478c-41d3-9bfa-405f5c0896c8',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Biaya Operasional Lapangan',
  21292500.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '46d86921-4c1b-46cd-b2cd-c91c55346757',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '8a19c98d-f320-4cb5-9f66-52c2a8bec4db',
  'Biaya Operasional Lapangan (FL)',
  14750000.00,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  'f6e646e9-92ad-41ea-a731-819d50abbeee',
  '6607332b-9424-4378-82ab-66f5b62ff95c',
  '0c59162d-be10-4517-9eeb-ebf9b9237b6d',
  'Kas & Bank - Kas Kecil Rizky',
  0.00,
  40870390.00,
  NOW(),
  NOW()
);

-- [JNL-JAN26-009] Penyusutan
-- Debit: 5,909,895.83  |  Kredit: 5,909,895.83
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  '3c087128-cc9b-443e-9d6e-f7761e122ef8',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-009',
  '2026-01-31',
  'Penyusutan',
  'general',
  'posted',
  5909895.83,
  5909895.83,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '6f6ecfdf-6dac-4077-b2b7-1ee768b503de',
  '3c087128-cc9b-443e-9d6e-f7761e122ef8',
  '173cc4cd-4f16-4e7c-b933-0ad00b150367',
  'Beban Penyusutan Kendaraan',
  5909895.83,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '6d297f32-5e7b-4a87-b538-8d4b8dd78a42',
  '3c087128-cc9b-443e-9d6e-f7761e122ef8',
  '6b941ecf-12ac-4948-8a63-6aedfa11e11c',
  'Akumulasi Penyusutan Kendaraan',
  0.00,
  5909895.83,
  NOW(),
  NOW()
);

-- [JNL-JAN26-010] Sewa Kantor Surabaya
-- Debit: 2,083,333.33  |  Kredit: 2,083,333.33
INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, date, description, type, status, total_debit, total_credit, created_at, updated_at)
VALUES (
  'ddf39737-3d93-461f-98be-554d0a51e7a4',
  'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
  '24520fa2-154d-4a24-b274-0d541367e6fe',
  'JNL-JAN26-010',
  '2026-01-31',
  'Sewa Kantor Surabaya',
  'general',
  'posted',
  2083333.33,
  2083333.33,
  NOW(),
  NOW()
);

INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '0f99845b-be77-4a18-b372-1a3f06e815b4',
  'ddf39737-3d93-461f-98be-554d0a51e7a4',
  '1a3bcab7-ee01-4cda-88ab-deb5919800a0',
  'Beban Sewa Gedung',
  2083333.33,
  0.00,
  NOW(),
  NOW()
);
INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
VALUES (
  '9ca1d69d-17a4-497a-8080-84300055f4be',
  'ddf39737-3d93-461f-98be-554d0a51e7a4',
  '44d11e40-8d2f-4e29-a93b-a6a7b510281f',
  'Sewa Gedung Dibayar Dimuka',
  0.00,
  2083333.33,
  NOW(),
  NOW()
);

COMMIT;