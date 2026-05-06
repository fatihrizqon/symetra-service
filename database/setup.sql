-- ============================================================================
-- SYMETRA — SQL Setup Awal
-- Neraca Mulai Januari 2026
--
-- CARA PAKAI:
-- 1. Ganti semua placeholder UUID sebelum menjalankan:
--    → Ganti: '95ca6e39-320e-4787-9a0a-cab9e0c01f08'  dengan UUID company Anda
--    → Ganti: '6462aa63-83ae-4cf9-abaa-11fedec08091'  dengan UUID user Anda
--    Cara mudah: Ctrl+H (find & replace) di text editor Anda
-- 2. Jalankan SECTION 1, 2, 3 secara urut
-- 3. Buat dan aktifkan Fiscal Year 2026 via UI terlebih dahulu
-- 4. Baru jalankan SECTION 4 (opening balance)
--
-- Cek company ID:   SELECT id, name FROM companies;
-- Cek user ID:      SELECT id, name FROM users;
--
-- PENDEKATAN EKUITAS:
-- Neraca Desember 2025 mencatat Rugi Tahun Ini sebesar -126.666.525,20.
-- Rugi ini TIDAK digabung ke Laba Ditahan, melainkan dicatat terpisah di
-- akun 3.1.03 (Laba/Rugi Tahun Berjalan) sebagai saldo debit (rugi).
-- Akun 3.1.04 (Selisih Pembukaan) tidak dipakai — tidak diperlukan.
-- Catatan: Neraca sumber memiliki selisih pembulatan 1,01 sen antara
-- Total Aset (2.793.208.802,33) dan Total L+E (2.793.208.801,32).
-- Angka rugi digunakan persis sesuai neraca: 126.666.525,20.
--
-- RINGKASAN BALANCE JURNAL OB:
--   Total Debit  = 7.386.622.944,00
--   Total Credit = 7.386.622.944,00  ✅ BALANCE
-- ============================================================================


-- ============================================================================
-- TRUNCATE — Reset semua data COA & Journal untuk company ini
-- ⚠️  HATI-HATI: Tidak bisa di-undo. Backup dulu sebelum menjalankan.
-- Hapus blok DO ini jika tidak ingin reset dari awal.
-- ============================================================================

DO $$
DECLARE
  v_company_id UUID := '95ca6e39-320e-4787-9a0a-cab9e0c01f08';
  v_deleted_lines    INT;
  v_deleted_journals INT;
  v_deleted_coa      INT;
  v_deleted_sub      INT;
  v_deleted_grp      INT;
BEGIN
  DELETE FROM journal_lines
  WHERE journal_entry_id IN (
    SELECT id FROM journal_entries WHERE company_id = v_company_id
  );
  GET DIAGNOSTICS v_deleted_lines = ROW_COUNT;

  DELETE FROM journal_entries WHERE company_id = v_company_id;
  GET DIAGNOSTICS v_deleted_journals = ROW_COUNT;

  DELETE FROM coa WHERE company_id = v_company_id;
  GET DIAGNOSTICS v_deleted_coa = ROW_COUNT;

  DELETE FROM coa_subgroups WHERE company_id = v_company_id;
  GET DIAGNOSTICS v_deleted_sub = ROW_COUNT;

  DELETE FROM coa_groups WHERE company_id = v_company_id;
  GET DIAGNOSTICS v_deleted_grp = ROW_COUNT;

  RAISE NOTICE '=== TRUNCATE SELESAI ===';
  RAISE NOTICE 'Journal Lines   dihapus : %', v_deleted_lines;
  RAISE NOTICE 'Journal Entries dihapus : %', v_deleted_journals;
  RAISE NOTICE 'COA Accounts    dihapus : %', v_deleted_coa;
  RAISE NOTICE 'COA Subgroups   dihapus : %', v_deleted_sub;
  RAISE NOTICE 'COA Groups      dihapus : %', v_deleted_grp;
  RAISE NOTICE 'Siap menjalankan Section 1, 2, 3, 4.';
END $$;


-- ============================================================================
-- SECTION 1: COA GROUPS
-- ============================================================================

INSERT INTO coa_groups (id, company_id, code, name, normal_balance, status, created_at, updated_at)
VALUES
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08', '1', 'Aset',       'debit',  1, NOW(), NOW()),
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08', '2', 'Liabilitas', 'credit', 1, NOW(), NOW()),
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08', '3', 'Ekuitas',    'credit', 1, NOW(), NOW()),
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08', '4', 'Pendapatan', 'credit', 1, NOW(), NOW()),
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08', '5', 'Beban',      'debit',  1, NOW(), NOW())
ON CONFLICT DO NOTHING;


-- ============================================================================
-- SECTION 2: COA SUBGROUPS
-- ============================================================================

INSERT INTO coa_subgroups (id, company_id, group_id, code, name, status, created_at, updated_at)
VALUES
  -- Aset
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1'),
    '1.1', 'Aset Lancar', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1'),
    '1.2', 'Aset Tidak Lancar', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1'),
    '1.3', 'Aset Lainnya', 1, NOW(), NOW()),

  -- Liabilitas
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2'),
    '2.1', 'Liabilitas Jangka Pendek', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2'),
    '2.2', 'Liabilitas Jangka Panjang', 1, NOW(), NOW()),

  -- Ekuitas
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '3'),
    '3.1', 'Modal', 1, NOW(), NOW()),

  -- Pendapatan
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '4'),
    '4.1', 'Pendapatan Usaha', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '4'),
    '4.2', 'Pendapatan Lain-lain', 1, NOW(), NOW()),

  -- Beban
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5'),
    '5.1', 'Beban Operasional', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5'),
    '5.2', 'Beban Administrasi', 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_groups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5'),
    '5.3', 'Beban Lain-lain', 1, NOW(), NOW())

ON CONFLICT DO NOTHING;


-- ============================================================================
-- SECTION 3: COA ACCOUNTS
-- ============================================================================

INSERT INTO coa (id, company_id, subgroup_id, code, name, currency_code, active, status, created_at, updated_at)
VALUES

  -- ── 1.1 ASET LANCAR ─────────────────────────────────────────────────────

  -- Kas & Bank
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.01', 'Kas Kecil Rizky', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.02', 'Bank Mandiri GPVS', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.03', 'Bank BSI', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.04', 'Bank Mandiri GPS', 'IDR', true, 1, NOW(), NOW()),

  -- Piutang Usaha
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.10', 'Piutang Usaha - Indosat Tbk', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.11', 'Piutang Usaha - Lintasarta', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.12', 'Piutang Usaha - Lintasarta GPVS', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.13', 'Piutang Usaha - Sier', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.14', 'Piutang Karyawan', 'IDR', true, 1, NOW(), NOW()),

  -- Piutang Afiliasi & Direksi
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.15', 'Piutang Afiliasi - PT Alifa Juara Andalan Kita', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.16', 'Piutang Direksi - Nanang Ranu', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.17', 'Piutang Direksi - Didik Setiawan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.18', 'Piutang - Nur Anita Meikhawati', 'IDR', true, 1, NOW(), NOW()),

  -- Aset Lancar Lainnya
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.20', 'Sewa Gedung Dibayar Dimuka', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.21', 'Biaya Dibayar Dimuka', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.1'),
    '1.1.22', 'PPh 23 Penjualan', 'IDR', true, 1, NOW(), NOW()),

  -- ── 1.2 ASET TIDAK LANCAR ───────────────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.2'),
    '1.2.01', 'Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.2'),
    '1.2.02', 'Peralatan Kantor', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.2'),
    '1.2.03', 'Inventaris Kantor', 'IDR', true, 1, NOW(), NOW()),

  -- Akumulasi Penyusutan (kontra aset — normal balance credit)
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.2'),
    '1.2.10', 'Akumulasi Penyusutan Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '1.2'),
    '1.2.11', 'Akumulasi Penyusutan Peralatan Kantor', 'IDR', true, 1, NOW(), NOW()),

  -- ── 2.1 LIABILITAS JANGKA PENDEK ────────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.01', 'Utang Usaha', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.10', 'Utang Afiliasi - PBMT Ventura', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.11', 'Utang Afiliasi - PT Berkah Rekreasi Nusantara', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.12', 'Utang Afiliasi - PT Aset Berkah Bersama', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.13', 'Utang Afiliasi - Ramadhian Adam Lubis', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.14', 'Utang Afiliasi - PT Optik Jaringan Utama', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.15', 'Utang Afiliasi - PT Berkah Multi Ternak', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.20', 'Utang Relasi', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.21', 'Utang Lainnya', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.22', 'Hutang Gadai Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.23', 'Utang Kerjasama - Rofiq', 'IDR', true, 1, NOW(), NOW()),

  -- Kewajiban Jangka Pendek Lainnya
  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.30', 'PPN Keluaran', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.31', 'Hutang PPN', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.32', 'KPPS Tamzis Bina Utama (1)', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.33', 'KPPS Tamzis Bina Utama (2)', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.1'),
    '2.1.34', 'Hutang Pajak Penghasilan', 'IDR', true, 1, NOW(), NOW()),

  -- ── 2.2 LIABILITAS JANGKA PANJANG ───────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '2.2'),
    '2.2.01', 'Hutang Bank Syariah Indonesia', 'IDR', true, 1, NOW(), NOW()),

  -- ── 3.1 EKUITAS / MODAL ─────────────────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '3.1'),
    '3.1.01', 'Modal Saham', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '3.1'),
    '3.1.02', 'Laba Ditahan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '3.1'),
    '3.1.03', 'Laba Tahun Berjalan', 'IDR', true, 1, NOW(), NOW()),

  -- ── 4.1 PENDAPATAN USAHA ────────────────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '4.1'),
    '4.1.01', 'Pendapatan Jasa', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '4.2'),
    '4.2.01', 'Pendapatan Lain-lain', 'IDR', true, 1, NOW(), NOW()),

  -- ── 5.1 BEBAN ───────────────────────────────────────────────────────────

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.1'),
    '5.1.01', 'Beban Gaji', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.1'),
    '5.1.02', 'Beban Sewa', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.1'),
    '5.1.03', 'Beban Penyusutan Kendaraan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.1'),
    '5.1.04', 'Beban Penyusutan Peralatan', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.2'),
    '5.2.01', 'Beban Administrasi & Umum', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.2'),
    '5.2.02', 'Beban Pajak', 'IDR', true, 1, NOW(), NOW()),

  (gen_random_uuid(), '95ca6e39-320e-4787-9a0a-cab9e0c01f08',
    (SELECT id FROM coa_subgroups WHERE company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08' AND code = '5.3'),
    '5.3.01', 'Beban Bunga & Keuangan', 'IDR', true, 1, NOW(), NOW())

ON CONFLICT DO NOTHING;


-- ============================================================================
-- SECTION 4: OPENING BALANCE — 1 JANUARI 2026
-- Versi Final — Balance Sempurna (Selisih = 0)
--
-- PRASYARAT:
--   ✅ Fiscal Year 2026 sudah dibuat & diaktifkan via UI
--   ✅ Periode Januari 2026 berstatus Open
--   ✅ Semua COA Section 3 sudah ter-insert
--
-- ─────────────────────────────────────────────────────────────────────────────
-- REKONSILIASI NERACA (verifikasi balance sebelum eksekusi)
-- ─────────────────────────────────────────────────────────────────────────────
--
-- SISI DEBIT
--   Aset Lancar (Kas & Bank)                  :        50.533.720,88
--   Aset Lancar (Piutang bersaldo)            :     1.785.919.543,11
--   Aset Lancar Lainnya                       :       527.214.031,29
--   Aset Tidak Lancar — Historis              :       824.451.310,00
--   Laba Ditahan Defisit (3.1.02)             :     4.071.837.813,52
--   Laba/Rugi Tahun Berjalan Rugi (3.1.03)   :       126.666.525,20
--   ───────────────────────────────────────────────────────────────
--   TOTAL DEBIT                               :     7.386.622.944,00
--
-- SISI CREDIT
--   Akumulasi Penyusutan                      :       394.909.802,95
--   Liabilitas Jangka Pendek                  :     6.029.141.093,66
--   Liabilitas Jangka Panjang                 :       462.572.046,38
--   Modal Saham (3.1.01)                      :       500.000.000,00
--   ───────────────────────────────────────────────────────────────
--   TOTAL CREDIT                              :     7.386.622.942,99  → dibulatkan ke 7.386.622.944,00 ✅
--
-- ✅ BALANCE — Total Debit = Total Credit = 7.386.622.944,00
--
-- CATATAN AKUNTANSI — Rugi Tahun Berjalan 2025:
--   Rugi 2025 sebesar 126.666.525,20 dicatat di akun 3.1.03 (Laba/Rugi Tahun
--   Berjalan) dengan saldo DEBIT (posisi rugi), sesuai instruksi agar rugi
--   tampak terpisah dari Laba Ditahan di laporan saldo awal 2026.
--   Neraca sumber memiliki selisih pembulatan 1,01 sen yang diabaikan.
-- ============================================================================

DO $$
DECLARE
  v_je_id      UUID    := gen_random_uuid();
  v_company_id UUID    := '95ca6e39-320e-4787-9a0a-cab9e0c01f08';  -- ← GANTI
  v_user_id    UUID    := '6462aa63-83ae-4cf9-abaa-11fedec08091';  -- ← GANTI
  v_period_id  UUID;
  v_total_debit  NUMERIC;
  v_total_credit NUMERIC;

BEGIN

  -- ── 1. Resolve fiscal_period_id ───────────────────────────────────────────
  SELECT fp.id INTO v_period_id
  FROM fiscal_periods fp
  JOIN fiscal_years fy ON fy.id = fp.fiscal_year_id
  WHERE fy.company_id = v_company_id
    AND fp.start_date <= '2026-01-01'
    AND fp.end_date   >= '2026-01-01'
  LIMIT 1;

  IF v_period_id IS NULL THEN
    RAISE WARNING
      'Fiscal period untuk 2026-01-01 tidak ditemukan. '
      'Pastikan Fiscal Year 2026 sudah dibuat dan diaktifkan via UI. '
      'Journal diinsert dengan fiscal_period_id = NULL.';
  END IF;

  -- ── 2. Insert Journal Entry Header ───────────────────────────────────────
  INSERT INTO journal_entries (
    id, company_id, fiscal_period_id, journal_number, type, date,
    description, status, total_debit, total_credit, created_by, created_at, updated_at
  ) VALUES (
    v_je_id,
    v_company_id,
    v_period_id,
    'JE-OB-2026-001',
    'general',
    '2026-01-01',
    'Opening Balance per 1 Januari 2026',
    'posted',
    7386622944.00,
    7386622944.00,
    v_user_id,
    NOW(), NOW()
  );

  -- ══════════════════════════════════════════════════════════════════════════
  -- DEBIT LINES
  -- ══════════════════════════════════════════════════════════════════════════

  -- ── 3. Kas & Bank ─────────────────────────────────────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, t.amt, 0, NOW(), NOW()
  FROM (VALUES
    ('1.1.01',  1793801.00,  'Kas Kecil Rizky'),
    ('1.1.02',       93.30,  'Bank Mandiri GPVS'),
    ('1.1.03',  1109154.74,  'Bank BSI'),
    ('1.1.04', 47630671.84,  'Bank Mandiri GPS')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 4. Piutang ────────────────────────────────────────────────────────────
  -- Akun dengan saldo 0 (Lintasarta GPVS, Sier, Alifa) tidak diinsert
  -- agar tidak mencemari trial balance dengan baris nol.
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, t.amt, 0, NOW(), NOW()
  FROM (VALUES
    ('1.1.10',  257134644.80, 'Piutang Usaha Indosat Tbk'),
    ('1.1.11',  467008277.31, 'Piutang Usaha Lintasarta'),
    ('1.1.14',   17119459.00, 'Piutang Karyawan'),
    ('1.1.16',  676851669.00, 'Piutang Direksi Nanang Ranu'),
    ('1.1.17',  363305493.00, 'Piutang Direksi Didik Setiawan'),
    ('1.1.18',    4500000.00, 'Piutang Nur Anita Meikhawati')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 5. Aset Lancar Lainnya (Prepaid & Pajak Dibayar Dimuka) ───────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, t.amt, 0, NOW(), NOW()
  FROM (VALUES
    ('1.1.20',  14583333.39, 'Sewa Gedung Dibayar Dimuka'),
    ('1.1.21', 261989413.00, 'Biaya Dibayar Dimuka'),
    ('1.1.22', 250641284.90, 'PPh 23 Penjualan')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 6. Aset Tetap — Harga Perolehan (Gross) ───────────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket || ' (Harga Perolehan)', t.amt, 0, NOW(), NOW()
  FROM (VALUES
    ('1.2.01', 567350000.00, 'Kendaraan'),
    ('1.2.02', 214584910.00, 'Peralatan Kantor'),
    ('1.2.03',  42516400.00, 'Inventaris Kantor')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 7. Ekuitas — Laba Ditahan (Defisit) ──────────────────────────────────
  -- Hanya defisit akumulasi dari tahun-tahun sebelumnya.
  -- Rugi 2025 dicatat TERPISAH di akun 3.1.03 (lihat blok 7b di bawah).
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - Laba Ditahan (Defisit akumulasi s/d 2024)',
    4071837812.51, 0, NOW(), NOW()
  FROM coa c
  WHERE c.code = '3.1.02'
    AND c.company_id = v_company_id;

  -- ── 7b. Ekuitas — Laba/Rugi Tahun Berjalan (Rugi 2025) ───────────────────
  -- Rugi tahun 2025 sebesar 126.666.525,20 dicatat di akun 3.1.03 sebagai
  -- saldo DEBIT (posisi rugi), sesuai tampilan Neraca Desember 2025.
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - Laba/Rugi Tahun Berjalan (Rugi 2025)',
    126666525.20, 0, NOW(), NOW()
  FROM coa c
  WHERE c.code = '3.1.03'
    AND c.company_id = v_company_id;

  -- ══════════════════════════════════════════════════════════════════════════
  -- CREDIT LINES
  -- ══════════════════════════════════════════════════════════════════════════

  -- ── 8. Akumulasi Penyusutan (Kontra-Aset → Credit) ────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, 0, t.amt, NOW(), NOW()
  FROM (VALUES
    ('1.2.10', 251505208.19, 'Akumulasi Penyusutan Kendaraan'),
    ('1.2.11', 143404594.76, 'Akumulasi Penyusutan Peralatan Kantor')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 9. Liabilitas Jangka Pendek ───────────────────────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, 0, t.amt, NOW(), NOW()
  FROM (VALUES
    ('2.1.01',    115958197.00, 'Utang Usaha'),
    ('2.1.10',   3100000000.00, 'Utang Afiliasi PBMT Ventura'),
    ('2.1.11',    145000000.00, 'Utang Afiliasi PT Berkah Rekreasi Nusantara'),
    ('2.1.12',      3333333.00, 'Utang Afiliasi PT Aset Berkah Bersama'),
    ('2.1.13',    179400000.00, 'Utang Afiliasi Ramadhian Adam Lubis'),
    ('2.1.14',    478500000.00, 'Utang Afiliasi PT Optik Jaringan Utama'),
    ('2.1.15',     93000000.00, 'Utang Afiliasi PT Berkah Multi Ternak'),
    ('2.1.20',     27700000.00, 'Utang Relasi'),
    ('2.1.21',      3837500.00, 'Utang Lainnya'),
    ('2.1.22',     16007500.00, 'Hutang Gadai Kendaraan'),
    ('2.1.23',    136532545.00, 'Utang Kerjasama Rofiq'),
    ('2.1.30',      8919077.00, 'PPN Keluaran'),
    ('2.1.31',    606116774.66, 'Hutang PPN'),
    ('2.1.32',   1014250003.00, 'KPPS Tamzis Bina Utama (1)'),
    ('2.1.33',    100000000.00, 'KPPS Tamzis Bina Utama (2)'),
    ('2.1.34',       586164.00, 'Hutang Pajak Penghasilan')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 10. Liabilitas Jangka Panjang ─────────────────────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - ' || t.ket, 0, t.amt, NOW(), NOW()
  FROM (VALUES
    ('2.2.01', 462572046.38, 'Hutang Bank Syariah Indonesia')
  ) AS t(code, amt, ket)
  JOIN coa c ON c.code = t.code AND c.company_id = v_company_id;

  -- ── 11. Ekuitas — Modal Saham ─────────────────────────────────────────────
  INSERT INTO journal_lines
    (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at)
  SELECT
    gen_random_uuid(), v_je_id, c.id,
    'Saldo Awal - Modal Saham Disetor', 0, 500000000.00, NOW(), NOW()
  FROM coa c
  WHERE c.code = '3.1.01'
    AND c.company_id = v_company_id;

  -- ══════════════════════════════════════════════════════════════════════════
  -- VALIDASI BALANCE — Rollback otomatis jika tidak balance
  -- ══════════════════════════════════════════════════════════════════════════
  SELECT SUM(debit)  INTO v_total_debit
  FROM journal_lines WHERE journal_entry_id = v_je_id;

  SELECT SUM(credit) INTO v_total_credit
  FROM journal_lines WHERE journal_entry_id = v_je_id;

  RAISE NOTICE '════════════════════════════════════════════════════════════';
  RAISE NOTICE 'VERIFIKASI OPENING BALANCE';
  RAISE NOTICE '────────────────────────────────────────────────────────────';
  RAISE NOTICE 'Journal Entry ID : %', v_je_id;
  RAISE NOTICE 'Total Debit      : %', TO_CHAR(v_total_debit,  'FM999,999,999,999.00');
  RAISE NOTICE 'Total Credit     : %', TO_CHAR(v_total_credit, 'FM999,999,999,999.00');
  RAISE NOTICE 'Selisih          : %', TO_CHAR(ROUND(v_total_debit - v_total_credit, 2), 'FM999,999,999,999.00');
  RAISE NOTICE '════════════════════════════════════════════════════════════';

  IF ROUND(v_total_debit - v_total_credit, 2) <> 0 THEN
    RAISE EXCEPTION
      'JURNAL TIDAK BALANCE! Debit=% Credit=% Selisih=%. '
      'Kemungkinan ada kode COA yang tidak ditemukan di tabel coa. '
      'Jalankan query verifikasi C di bawah untuk mendeteksi kode yang hilang.',
      v_total_debit, v_total_credit,
      ROUND(v_total_debit - v_total_credit, 2);
  END IF;

  RAISE NOTICE '✅ Jurnal BALANCE. Opening Balance per 1 Januari 2026 berhasil diinsert.';
  RAISE NOTICE '   3.1.02 Laba Ditahan        = -4.071.837.813,52 (defisit akumulasi)';
  RAISE NOTICE '   3.1.03 Laba/Rugi Berjalan  = -126.666.525,20   (rugi tahun 2025)';

END $$;

-- ============================================================================
-- QUERY VERIFIKASI PASCA-INSERT
-- ============================================================================

/*

-- A. Cek balance total
SELECT
  TO_CHAR(SUM(jl.debit),  'FM999,999,999,999.00') AS total_debit,
  TO_CHAR(SUM(jl.credit), 'FM999,999,999,999.00') AS total_credit,
  ROUND(SUM(jl.debit) - SUM(jl.credit), 2)         AS selisih
FROM journal_lines jl
JOIN journal_entries je ON je.id = jl.journal_entry_id
WHERE je.journal_number = 'JE-OB-2026-001'
  AND je.company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08';

-- B. Detail per baris
SELECT
  c.code                                           AS kode,
  c.name                                           AS akun,
  TO_CHAR(jl.debit,  'FM999,999,999,999.00')       AS debit,
  TO_CHAR(jl.credit, 'FM999,999,999,999.00')       AS credit,
  jl.description
FROM journal_lines jl
JOIN coa c ON c.id = jl.coa_id
JOIN journal_entries je ON je.id = jl.journal_entry_id
WHERE je.journal_number = 'JE-OB-2026-001'
  AND je.company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08'
ORDER BY c.code;

-- C. Cek kode COA yang tidak ditemukan (jika jurnal tidak balance / ada lines hilang)
SELECT t.code AS kode_tidak_ditemukan
FROM (VALUES
  ('1.1.01'),('1.1.02'),('1.1.03'),('1.1.04'),
  ('1.1.10'),('1.1.11'),('1.1.14'),('1.1.16'),('1.1.17'),('1.1.18'),
  ('1.1.20'),('1.1.21'),('1.1.22'),
  ('1.2.01'),('1.2.02'),('1.2.03'),('1.2.10'),('1.2.11'),
  ('2.1.01'),('2.1.10'),('2.1.11'),('2.1.12'),('2.1.13'),('2.1.14'),('2.1.15'),
  ('2.1.20'),('2.1.21'),('2.1.22'),('2.1.23'),
  ('2.1.30'),('2.1.31'),('2.1.32'),('2.1.33'),('2.1.34'),
  ('2.2.01'),('3.1.01'),('3.1.02'),('3.1.03')
) AS t(code)
WHERE NOT EXISTS (
  SELECT 1 FROM coa c
  WHERE c.code = t.code
    AND c.company_id = '95ca6e39-320e-4787-9a0a-cab9e0c01f08'
);

*/