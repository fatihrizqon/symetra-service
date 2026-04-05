-- ============================================================================
-- SYMETRA — SQL Setup Awal
-- Neraca Mulai Januari 2026
--
-- CARA PAKAI:
-- 1. Ganti semua :company_id dengan UUID company Anda
--    → Cari di tabel companies: SELECT id, name FROM companies;
-- 2. Ganti :user_id dengan UUID user yang melakukan setup
--    → Cari di tabel users: SELECT id, name FROM users;
-- 3. Jalankan SECTION 1, 2, 3 secara urut
-- 4. Setelah SECTION 3, jalankan SECTION 4 (opening balance) SETELAH
--    Fiscal Year 2026 dibuat dan diaktifkan via UI
-- ============================================================================

-- ============================================================================
-- GANTI DUA VARIABLE INI SEBELUM MENJALANKAN
-- ============================================================================
-- Cara mudah: Ctrl+H (find & replace) di text editor Anda
--   Ganti: 'f90772d4-4be7-43cd-a649-0c4bc7a3d953'  → UUID company Anda
--   Ganti: '509c15ac-8a5a-4572-b343-ac46a7251279'     → UUID user Anda

-- Cek company ID:   SELECT id, name FROM companies;
-- Cek user ID:      SELECT id, name FROM users;

-- ============================================================================
-- SECTION 1: COA GROUPS
-- Hierarki Akuntansi Standard Indonesia
-- normal_balance: 'debit' untuk Aset & Beban, 'credit' untuk Liabilitas, Ekuitas, Pendapatan
-- ============================================================================

INSERT INTO coa_groups (id, company_id, code, name, normal_balance, status, created_at, updated_at)
VALUES
  -- Kelompok Aset
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953', '1', 'Aset',        'debit',  1, NOW(), NOW()),
  -- Kelompok Liabilitas
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953', '2', 'Liabilitas',  'credit', 1, NOW(), NOW()),
  -- Kelompok Ekuitas
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953', '3', 'Ekuitas',     'credit', 1, NOW(), NOW()),
  -- Kelompok Pendapatan
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953', '4', 'Pendapatan',  'credit', 1, NOW(), NOW()),
  -- Kelompok Beban
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953', '5', 'Beban',       'debit',  1, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SECTION 2: COA SUBGROUPS
-- Menggunakan subquery untuk ambil group_id dari nama group
-- ============================================================================

INSERT INTO coa_subgroups (id, company_id, group_id, code, name, status, created_at, updated_at)
VALUES
  -- ── Subgroup Aset ──────────────────────────────────────────────────────
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1'),
    '1.1', 'Aset Lancar', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1'),
    '1.2', 'Aset Tidak Lancar', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1'),
    '1.3', 'Aset Lainnya', 1, NOW(), NOW()),

  -- ── Subgroup Liabilitas ─────────────────────────────────────────────────
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2'),
    '2.1', 'Liabilitas Jangka Pendek', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2'),
    '2.2', 'Liabilitas Jangka Panjang', 1, NOW(), NOW()),

  -- ── Subgroup Ekuitas ────────────────────────────────────────────────────
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '3'),
    '3.1', 'Modal', 1, NOW(), NOW()),

  -- ── Subgroup Pendapatan ─────────────────────────────────────────────────
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '4'),
    '4.1', 'Pendapatan Usaha', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '4'),
    '4.2', 'Pendapatan Lain-lain', 1, NOW(), NOW()),

  -- ── Subgroup Beban ──────────────────────────────────────────────────────
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5'),
    '5.1', 'Beban Operasional', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5'),
    '5.2', 'Beban Administrasi', 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_groups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5'),
    '5.3', 'Beban Lain-lain', 1, NOW(), NOW())

ON CONFLICT DO NOTHING;

-- ============================================================================
-- SECTION 3: COA ACCOUNTS
-- Berdasarkan akun-akun yang muncul di Neraca Mulai Januari 2026
-- ============================================================================

INSERT INTO coa (id, company_id, subgroup_id, code, name, currency_code, active, status, created_at, updated_at)
VALUES

  -- ═══════════════════════════════════════════════════════════════════════
  -- 1.1 ASET LANCAR
  -- ═══════════════════════════════════════════════════════════════════════

  -- Kas dan Setara Kas
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.01', 'Kas Kecil Rizky', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.02', 'Bank Mandiri GPVS', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.03', 'Bank BSI', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.04', 'Bank Mandiri GPS', 'IDR', true, 1, NOW(), NOW()),

  -- Piutang Usaha
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.10', 'Piutang Usaha - Indosat Tbk', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.11', 'Piutang Usaha - Lintasarta', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.12', 'Piutang Usaha - Lintasarta GPVS', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.13', 'Piutang Usaha - Sier', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.14', 'Piutang Karyawan', 'IDR', true, 1, NOW(), NOW()),

  -- Piutang Afiliasi & Direksi
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.15', 'Piutang Afiliasi - PT Alifa Juara Andalan Kita', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.16', 'Piutang Direksi - Nanang Ranu', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.17', 'Piutang Direksi - Didik Setiawan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.18', 'Piutang - Nur Anita Meikhawati', 'IDR', true, 1, NOW(), NOW()),

  -- Aset Lancar Lainnya
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.20', 'Sewa Gedung Dibayar Dimuka', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.21', 'Biaya Dibayar Dimuka', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.1'),
    '1.1.22', 'PPh 23 Penjualan', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 1.2 ASET TIDAK LANCAR
  -- ═══════════════════════════════════════════════════════════════════════

  -- Nilai Historis Aset Tetap
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.2'),
    '1.2.01', 'Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.2'),
    '1.2.02', 'Peralatan Kantor', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.2'),
    '1.2.03', 'Inventaris Kantor', 'IDR', true, 1, NOW(), NOW()),

  -- Akumulasi Penyusutan (Kontra Aset — normal balance credit, tapi di grup Aset)
  -- Di-input sebagai nilai negatif di debit (atau positif di credit) saat opening balance
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.2'),
    '1.2.10', 'Akumulasi Penyusutan Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '1.2'),
    '1.2.11', 'Akumulasi Penyusutan Peralatan Kantor', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 2.1 LIABILITAS JANGKA PENDEK
  -- ═══════════════════════════════════════════════════════════════════════

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.01', 'Utang Usaha', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.10', 'Utang Afiliasi - PBMT Ventura', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.11', 'Utang Afiliasi - PT Berkah Rekreasi Nusantara', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.12', 'Utang Afiliasi - PT Aset Berkah Bersama', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.13', 'Utang Afiliasi - Ramadhian Adam Lubis', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.14', 'Utang Afiliasi - PT Optik Jaringan Utama', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.15', 'Utang Afiliasi - PT Berkah Multi Ternak', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.20', 'Utang Relasi', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.21', 'Utang Lainnya', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.22', 'Hutang Gadai Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.23', 'Utang Kerjasama - Rofiq', 'IDR', true, 1, NOW(), NOW()),

  -- Kewajiban Jangka Pendek Lainnya (Pajak)
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.30', 'PPN Keluaran', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.31', 'Hutang PPN', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.32', 'KPPS Tamzis Bina Utama (1)', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.33', 'KPPS Tamzis Bina Utama (2)', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.1'),
    '2.1.34', 'Hutang Pajak Penghasilan', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 2.2 LIABILITAS JANGKA PANJANG
  -- ═══════════════════════════════════════════════════════════════════════

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '2.2'),
    '2.2.01', 'Hutang Bank Syariah Indonesia', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 3.1 EKUITAS / MODAL
  -- ═══════════════════════════════════════════════════════════════════════

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '3.1'),
    '3.1.01', 'Modal Saham', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '3.1'),
    '3.1.02', 'Laba Ditahan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '3.1'),
    '3.1.03', 'Laba Tahun Berjalan', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 4.1 PENDAPATAN USAHA (template untuk operasional)
  -- ═══════════════════════════════════════════════════════════════════════

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '4.1'),
    '4.1.01', 'Pendapatan Jasa', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '4.2'),
    '4.2.01', 'Pendapatan Lain-lain', 'IDR', true, 1, NOW(), NOW()),

  -- ═══════════════════════════════════════════════════════════════════════
  -- 5.1 BEBAN OPERASIONAL (template — tambahkan sesuai kebutuhan)
  -- ═══════════════════════════════════════════════════════════════════════

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.1'),
    '5.1.01', 'Beban Gaji', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.1'),
    '5.1.02', 'Beban Sewa', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.1'),
    '5.1.03', 'Beban Penyusutan Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.1'),
    '5.1.04', 'Beban Penyusutan Peralatan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.2'),
    '5.2.01', 'Beban Administrasi & Umum', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.2'),
    '5.2.02', 'Beban Pajak', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '5.3'),
    '5.3.01', 'Beban Bunga & Keuangan', 'IDR', true, 1, NOW(), NOW()),

  -- Akun sementara untuk selisih pembukaan (jika neraca tidak balance)
  (gen_random_uuid(), 'f90772d4-4be7-43cd-a649-0c4bc7a3d953',
    (SELECT id FROM coa_subgroups WHERE company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953' AND code = '3.1'),
    '3.1.99', 'Selisih Pembukaan (Sementara)', 'IDR', true, 1, NOW(), NOW())

ON CONFLICT DO NOTHING;

-- ============================================================================
-- VERIFIKASI: Cek hasil insert
-- ============================================================================

-- SELECT g.code AS grp, g.name AS group_name, s.code AS sub, s.name AS subgroup_name,
--        c.code AS coa_code, c.name AS coa_name
-- FROM coa c
-- JOIN coa_subgroups s ON s.id = c.subgroup_id
-- JOIN coa_groups g ON g.id = s.group_id
-- WHERE c.company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953'
-- ORDER BY c.code;

-- ============================================================================
-- SECTION 4: OPENING BALANCE — JURNAL SALDO AWAL 31 DESEMBER 2026
--
-- PRASYARAT SEBELUM MENJALANKAN INI:
-- ✅ Fiscal Year 2026 sudah dibuat via UI (2026-01-01 s/d 2026-12-31)
-- ✅ Fiscal Year 2026 sudah di-Activate
-- ✅ Periode Januari 2026 berstatus Open
-- ✅ Semua COA di Section 3 sudah ter-insert
--
-- CATATAN NERACA:
-- Total Aset      = 2.793.208.802,33
-- Total Liabilitas= 6.491.713.140,04
-- Modal Saham     = 500.000.000,00
-- Laba Ditahan    = -4.071.837.813,52 (DEFISIT → posisi DEBIT)
-- Laba Tahun Ini  = 0
--
-- Cek balance:
--   DEBIT  : Aset (2.793.208.802,33) + Laba Ditahan/Defisit (4.071.837.813,52)
--            + Selisih (126.666.524,19) = 6.991.713.139,04 (pembulatan)
--   CREDIT : Liabilitas (6.491.713.140,04) + Modal Saham (500.000.000)
--            = 6.991.713.140,04
--
-- Selisih neraca: JUMLAH ASET (2.793.208.802,33) ≠ JUMLAH L+E (2.919.875.326,52)
-- Selisih = 2.919.875.326,52 - 2.793.208.802,33 = 126.666.524,19
-- → Dimasukkan ke akun "Selisih Pembukaan (Sementara)" di sisi DEBIT
--   untuk menyeimbangkan jurnal.
-- → Investigasi selisih ini sebelum menutup periode pertama!
-- ============================================================================

DO $$
DECLARE
  v_je_id        UUID := gen_random_uuid();
  v_company_id   UUID := 'f90772d4-4be7-43cd-a649-0c4bc7a3d953';
  v_user_id      UUID := '509c15ac-8a5a-4572-b343-ac46a7251279';

  -- Ambil fiscal_period_id untuk Desember 2026
  -- Jika Fiscal Year 2026 tidak ada, gunakan periode Januari 2027
  -- (Opening balance di tanggal 2026-12-31 akan fall ke period Desember 2026 jika ada)
  v_period_id    UUID;

  -- COA IDs — diambil dinamis berdasarkan kode akun
  coa_kas_kecil_rizky     UUID;
  coa_bank_mandiri_gpvs   UUID;
  coa_bank_bsi            UUID;
  coa_bank_mandiri_gps    UUID;
  coa_piutang_indosat     UUID;
  coa_piutang_lintasarta  UUID;
  coa_piutang_lintasarta2 UUID;
  coa_piutang_sier        UUID;
  coa_piutang_karyawan    UUID;
  coa_piutang_afiliasi    UUID;
  coa_piutang_nanang      UUID;
  coa_piutang_didik       UUID;
  coa_piutang_nuranita    UUID;
  coa_sewa_dimuka         UUID;
  coa_biaya_dimuka        UUID;
  coa_pph23               UUID;
  coa_kendaraan           UUID;
  coa_peralatan           UUID;
  coa_inventaris          UUID;
  coa_akum_kendaraan      UUID;
  coa_akum_peralatan      UUID;
  coa_utang_usaha         UUID;
  coa_utang_pbmt          UUID;
  coa_utang_berkah_rek    UUID;
  coa_utang_aset_berkah   UUID;
  coa_utang_ramadhian     UUID;
  coa_utang_optik         UUID;
  coa_utang_ternak        UUID;
  coa_utang_relasi        UUID;
  coa_utang_lainnya       UUID;
  coa_gadai_kendaraan     UUID;
  coa_utang_rofiq         UUID;
  coa_ppn_keluaran        UUID;
  coa_hutang_ppn          UUID;
  coa_kpps_1              UUID;
  coa_kpps_2              UUID;
  coa_hutang_pph          UUID;
  coa_hutang_bsi_panjang  UUID;
  coa_modal_saham         UUID;
  coa_laba_ditahan        UUID;
  coa_selisih             UUID;

BEGIN

  -- ── Ambil fiscal_period_id ──────────────────────────────────────────────
  -- Cari period yang mencakup tanggal 2026-1-1 (FY 2026)
  -- Jika tidak ada (FY 2026 belum dibuat), set NULL
  SELECT fp.id INTO v_period_id
  FROM fiscal_periods fp
  JOIN fiscal_years fy ON fy.id = fp.fiscal_year_id
  WHERE fy.company_id = v_company_id
    AND fp.start_date <= '2026-1-1'
    AND fp.end_date   >= '2026-1-1'
  LIMIT 1;
  -- Jika NULL, journal tetap bisa dibuat (fiscal_period_id nullable)

  -- ── Ambil semua COA ID ──────────────────────────────────────────────────
  SELECT id INTO coa_kas_kecil_rizky    FROM coa WHERE company_id = v_company_id AND code = '1.1.01';
  SELECT id INTO coa_bank_mandiri_gpvs  FROM coa WHERE company_id = v_company_id AND code = '1.1.02';
  SELECT id INTO coa_bank_bsi           FROM coa WHERE company_id = v_company_id AND code = '1.1.03';
  SELECT id INTO coa_bank_mandiri_gps   FROM coa WHERE company_id = v_company_id AND code = '1.1.04';
  SELECT id INTO coa_piutang_indosat    FROM coa WHERE company_id = v_company_id AND code = '1.1.10';
  SELECT id INTO coa_piutang_lintasarta FROM coa WHERE company_id = v_company_id AND code = '1.1.11';
  SELECT id INTO coa_piutang_lintasarta2 FROM coa WHERE company_id = v_company_id AND code = '1.1.12';
  SELECT id INTO coa_piutang_sier       FROM coa WHERE company_id = v_company_id AND code = '1.1.13';
  SELECT id INTO coa_piutang_karyawan   FROM coa WHERE company_id = v_company_id AND code = '1.1.14';
  SELECT id INTO coa_piutang_afiliasi   FROM coa WHERE company_id = v_company_id AND code = '1.1.15';
  SELECT id INTO coa_piutang_nanang     FROM coa WHERE company_id = v_company_id AND code = '1.1.16';
  SELECT id INTO coa_piutang_didik      FROM coa WHERE company_id = v_company_id AND code = '1.1.17';
  SELECT id INTO coa_piutang_nuranita   FROM coa WHERE company_id = v_company_id AND code = '1.1.18';
  SELECT id INTO coa_sewa_dimuka        FROM coa WHERE company_id = v_company_id AND code = '1.1.20';
  SELECT id INTO coa_biaya_dimuka       FROM coa WHERE company_id = v_company_id AND code = '1.1.21';
  SELECT id INTO coa_pph23              FROM coa WHERE company_id = v_company_id AND code = '1.1.22';
  SELECT id INTO coa_kendaraan          FROM coa WHERE company_id = v_company_id AND code = '1.2.01';
  SELECT id INTO coa_peralatan          FROM coa WHERE company_id = v_company_id AND code = '1.2.02';
  SELECT id INTO coa_inventaris         FROM coa WHERE company_id = v_company_id AND code = '1.2.03';
  SELECT id INTO coa_akum_kendaraan     FROM coa WHERE company_id = v_company_id AND code = '1.2.10';
  SELECT id INTO coa_akum_peralatan     FROM coa WHERE company_id = v_company_id AND code = '1.2.11';
  SELECT id INTO coa_utang_usaha        FROM coa WHERE company_id = v_company_id AND code = '2.1.01';
  SELECT id INTO coa_utang_pbmt         FROM coa WHERE company_id = v_company_id AND code = '2.1.10';
  SELECT id INTO coa_utang_berkah_rek   FROM coa WHERE company_id = v_company_id AND code = '2.1.11';
  SELECT id INTO coa_utang_aset_berkah  FROM coa WHERE company_id = v_company_id AND code = '2.1.12';
  SELECT id INTO coa_utang_ramadhian    FROM coa WHERE company_id = v_company_id AND code = '2.1.13';
  SELECT id INTO coa_utang_optik        FROM coa WHERE company_id = v_company_id AND code = '2.1.14';
  SELECT id INTO coa_utang_ternak       FROM coa WHERE company_id = v_company_id AND code = '2.1.15';
  SELECT id INTO coa_utang_relasi       FROM coa WHERE company_id = v_company_id AND code = '2.1.20';
  SELECT id INTO coa_utang_lainnya      FROM coa WHERE company_id = v_company_id AND code = '2.1.21';
  SELECT id INTO coa_gadai_kendaraan    FROM coa WHERE company_id = v_company_id AND code = '2.1.22';
  SELECT id INTO coa_utang_rofiq        FROM coa WHERE company_id = v_company_id AND code = '2.1.23';
  SELECT id INTO coa_ppn_keluaran       FROM coa WHERE company_id = v_company_id AND code = '2.1.30';
  SELECT id INTO coa_hutang_ppn         FROM coa WHERE company_id = v_company_id AND code = '2.1.31';
  SELECT id INTO coa_kpps_1             FROM coa WHERE company_id = v_company_id AND code = '2.1.32';
  SELECT id INTO coa_kpps_2             FROM coa WHERE company_id = v_company_id AND code = '2.1.33';
  SELECT id INTO coa_hutang_pph         FROM coa WHERE company_id = v_company_id AND code = '2.1.34';
  SELECT id INTO coa_hutang_bsi_panjang FROM coa WHERE company_id = v_company_id AND code = '2.2.01';
  SELECT id INTO coa_modal_saham        FROM coa WHERE company_id = v_company_id AND code = '3.1.01';
  SELECT id INTO coa_laba_ditahan       FROM coa WHERE company_id = v_company_id AND code = '3.1.02';
  SELECT id INTO coa_selisih            FROM coa WHERE company_id = v_company_id AND code = '3.1.99';

  -- ── Insert Journal Entry Header ─────────────────────────────────────────
  INSERT INTO journal_entries (
    id, company_id, fiscal_period_id, journal_number, type, date,
    description, status, total_debit, total_credit, created_by, created_at, updated_at
  ) VALUES (
    v_je_id,
    v_company_id,
    v_period_id,
    'JE-OB-2026-001',
    'general',
    '2026-1-1',
    'Saldo Awal per 1 Januari 2026 (Opening Balance dari Neraca Desember 2025)',
    'posted',
    -- Total Debit: Aset 2.793.208.802,33 + Defisit Laba Ditahan 4.071.837.813,52 + Selisih 126.666.524,19
    6991713139.04,
    -- Total Credit: Liabilitas 6.491.713.140,04 + Modal Saham 500.000.000
    6991713140.04,
    v_user_id,
    NOW(), NOW()
  );

  -- ── Insert Journal Lines ────────────────────────────────────────────────

  -- ════════════════════════════════════════
  -- DEBIT: ASET LANCAR — Kas & Bank
  -- ════════════════════════════════════════
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
    (gen_random_uuid(), v_je_id, coa_kas_kecil_rizky,
      'Saldo Awal - Kas Kecil Rizky', 1793801.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_bank_mandiri_gpvs,
      'Saldo Awal - Bank Mandiri GPVS', 93.30, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_bank_bsi,
      'Saldo Awal - Bank BSI', 1109154.74, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_bank_mandiri_gps,
      'Saldo Awal - Bank Mandiri GPS', 47630671.84, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- DEBIT: ASET LANCAR — Piutang
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_piutang_indosat,
      'Saldo Awal - Piutang Indosat Tbk', 257134644.80, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_lintasarta,
      'Saldo Awal - Piutang Lintasarta', 467008277.31, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_lintasarta2,
      'Saldo Awal - Piutang Lintasarta GPVS', 0.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_sier,
      'Saldo Awal - Piutang Sier', 0.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_karyawan,
      'Saldo Awal - Piutang Karyawan', 17119459.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_afiliasi,
      'Saldo Awal - Piutang Afiliasi PT Alifa', 0.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_nanang,
      'Saldo Awal - Piutang Direksi Nanang Ranu', 676851669.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_didik,
      'Saldo Awal - Piutang Direksi Didik Setiawan', 363305493.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_piutang_nuranita,
      'Saldo Awal - Piutang Nur Anita Meikhawati', 4500000.00, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- DEBIT: ASET LANCAR — Aset Lainnya
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_sewa_dimuka,
      'Saldo Awal - Sewa Gedung Dibayar Dimuka', 14583333.39, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_biaya_dimuka,
      'Saldo Awal - Biaya Dibayar Dimuka', 261989413.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_pph23,
      'Saldo Awal - PPh 23 Penjualan', 250641284.90, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- DEBIT: ASET TIDAK LANCAR — Nilai Historis
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_kendaraan,
      'Saldo Awal - Kendaraan (Harga Perolehan)', 567350000.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_peralatan,
      'Saldo Awal - Peralatan Kantor (Harga Perolehan)', 214584910.00, 0, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_inventaris,
      'Saldo Awal - Inventaris Kantor', 42516400.00, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- CREDIT: ASET TIDAK LANCAR — Akumulasi Penyusutan
  -- (Kontra aset: normal balance CREDIT)
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_akum_kendaraan,
      'Saldo Awal - Akum. Penyusutan Kendaraan', 0, 251505208.19, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_akum_peralatan,
      'Saldo Awal - Akum. Penyusutan Peralatan', 0, 143404594.76, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- DEBIT: EKUITAS — Laba Ditahan DEFISIT
  -- Laba Ditahan = -4.071.837.813,52 (negatif = defisit = posisi DEBIT)
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_laba_ditahan,
      'Saldo Awal - Defisit Laba Ditahan', 4071837813.52, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- DEBIT: Selisih Pembukaan (Sementara)
  -- Selisih antara Total Aset dengan Total L+E di neraca
  -- Investigasi dan eliminasi sebelum tutup periode pertama
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_selisih,
      'Selisih Pembukaan - Investigasi lebih lanjut', 126666524.19, 0, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- CREDIT: LIABILITAS JANGKA PENDEK — Utang Usaha
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_utang_usaha,
      'Saldo Awal - Utang Usaha', 0, 115958197.00, NOW(), NOW()),

  -- CREDIT: Utang Afiliasi
    (gen_random_uuid(), v_je_id, coa_utang_pbmt,
      'Saldo Awal - Utang Afiliasi PBMT Ventura', 0, 3100000000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_berkah_rek,
      'Saldo Awal - Utang Afiliasi PT Berkah Rekreasi Nusantara', 0, 145000000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_aset_berkah,
      'Saldo Awal - Utang Afiliasi PT Aset Berkah Bersama', 0, 3333333.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_ramadhian,
      'Saldo Awal - Utang Afiliasi Ramadhian Adam Lubis', 0, 179400000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_optik,
      'Saldo Awal - Utang Afiliasi PT Optik Jaringan Utama', 0, 478500000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_ternak,
      'Saldo Awal - Utang Afiliasi PT Berkah Multi Ternak', 0, 93000000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_relasi,
      'Saldo Awal - Utang Relasi', 0, 27700000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_lainnya,
      'Saldo Awal - Utang Lainnya', 0, 3837500.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_gadai_kendaraan,
      'Saldo Awal - Hutang Gadai Kendaraan', 0, 16007500.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_utang_rofiq,
      'Saldo Awal - Utang Kerjasama Rofiq', 0, 136532545.00, NOW(), NOW()),

  -- CREDIT: Kewajiban Jangka Pendek Lainnya
    (gen_random_uuid(), v_je_id, coa_ppn_keluaran,
      'Saldo Awal - PPN Keluaran', 0, 8919077.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_hutang_ppn,
      'Saldo Awal - Hutang PPN', 0, 606116774.66, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_kpps_1,
      'Saldo Awal - KPPS Tamzis Bina Utama (1)', 0, 1014250003.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_kpps_2,
      'Saldo Awal - KPPS Tamzis Bina Utama (2)', 0, 100000000.00, NOW(), NOW()),

    (gen_random_uuid(), v_je_id, coa_hutang_pph,
      'Saldo Awal - Hutang Pajak Penghasilan', 0, 586164.00, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- CREDIT: LIABILITAS JANGKA PANJANG
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_hutang_bsi_panjang,
      'Saldo Awal - Hutang Bank Syariah Indonesia', 0, 462572046.38, NOW(), NOW()),

  -- ════════════════════════════════════════
  -- CREDIT: EKUITAS — Modal Saham
  -- ════════════════════════════════════════
    (gen_random_uuid(), v_je_id, coa_modal_saham,
      'Saldo Awal - Modal Saham', 0, 500000000.00, NOW(), NOW());

  RAISE NOTICE 'Opening Balance berhasil diinsert. Journal ID: %', v_je_id;
  RAISE NOTICE 'Cek di UI: Journal Entries → filter status=posted → JE-OB-2026-001';
  RAISE NOTICE '⚠️  Selisih 126.666.524,19 di akun 3.1.99 perlu diinvestigasi!';

END $$;

-- ============================================================================
-- VERIFIKASI OPENING BALANCE
-- Jalankan query ini setelah Section 4 untuk memastikan jurnal balance
-- ============================================================================

/*
SELECT
  jl.description,
  c.code,
  c.name AS akun,
  jl.debit,
  jl.credit
FROM journal_lines jl
JOIN coa c ON c.id = jl.coa_id
JOIN journal_entries je ON je.id = jl.journal_entry_id
WHERE je.journal_number = 'JE-OB-2026-001'
  AND je.company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953'
ORDER BY c.code;

-- Cek balance:
SELECT
  SUM(debit)  AS total_debit,
  SUM(credit) AS total_credit,
  SUM(debit) - SUM(credit) AS selisih
FROM journal_lines jl
JOIN journal_entries je ON je.id = jl.journal_entry_id
WHERE je.journal_number = 'JE-OB-2026-001'
  AND je.company_id = 'f90772d4-4be7-43cd-a649-0c4bc7a3d953';
*/

-- ============================================================================
-- LANGKAH SELANJUTNYA SETELAH SQL INI:
-- 1. Buka Symetra UI → Reports → Trial Balance (per 2026-1-1)
--    Bandingkan dengan neraca yang ada
-- 2. Investigasi akun 3.1.99 (Selisih Pembukaan)
--    Temukan akun yang salah/terlewat dari neraca asli
-- 3. Setelah selisih = 0, hapus atau zeroing akun 3.1.99
-- 4. Mulai input transaksi Januari 2026
-- ============================================================================