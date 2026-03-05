-- =============================================================
-- SEED: Journal Entries — Penjualan Telur September 2025
-- Strategi: 1 entry per hari per metode pembayaran
-- Setiap entry: 4 lines (Bank/Cash debit, HPP debit, Sales credit, Inventory credit)
--
-- PENTING: Jalankan query berikut dulu untuk dapat coa IDs kamu:
-- SELECT id, code, name FROM coa WHERE code IN
--   ('10101','10102','10104','40101','50101');
-- =============================================================

DO $$
DECLARE
  -- coa IDs — ambil dari database kamu
  v_cash_id    UUID := (SELECT id FROM coa WHERE code = '10101');  -- Cash
  v_bank_id    UUID := (SELECT id FROM coa WHERE code = '10102');  -- Bank
  v_inventory  UUID := (SELECT id FROM coa WHERE code = '10104');  -- Inventory - Eggs
  v_sales      UUID := (SELECT id FROM coa WHERE code = '40101');  -- Egg Sales
  v_cogs       UUID := (SELECT id FROM coa WHERE code = '50101');  -- Egg Purchase Cost
  v_admin_id   UUID := (SELECT id FROM users WHERE email = 'admin@example.com' LIMIT 1);
  v_je_id      UUID;
BEGIN

  -- ── JE-2509-0001: 2025-09-10 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0001',
    '2025-09-10',
    'Penjualan Telur Bebek — Dyah - Kadipiro, Astrini - Mlati, Ida - SMPN 3 YK [Bank BPD]',
    'posted',
    238500.00,   -- total debit (bank + cogs)
    238500.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  135000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   103500.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 135000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 103500.00);

  -- ── JE-2509-0002: 2025-09-12 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0002',
    '2025-09-12',
    'Penjualan Telur Bebek — Shanti Bintaran, Nila - DISPERTARU, Qayyim - DISPERTARU, +2 lainnya [Bank BPD]',
    'posted',
    344500.00,   -- total debit (bank + cogs)
    344500.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  195000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   149500.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 195000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 149500.00);

  -- ── JE-2509-0003: 2025-09-16 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0003',
    '2025-09-16',
    'Penjualan Telur Ayam, Telur Asin — Shanti SMK, Aris - DISPERTARU, Astri - DISPERTARU, +2 lainnya [Bank BPD]',
    'posted',
    324830.00,   -- total debit (bank + cogs)
    324830.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  187500.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   137330.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 187500.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 137330.00);

  -- ── JE-2509-0004: 2025-09-18 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0004',
    '2025-09-18',
    'Penjualan Telur Ayam — Ida - SMPN 3 YK, Ayu DISPERTARU [Bank BPD]',
    'posted',
    124140.00,   -- total debit (bank + cogs)
    124140.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  75000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   49140.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 75000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 49140.00);

  -- ── JE-2509-0005: 2025-09-19 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0005',
    '2025-09-19',
    'Penjualan Telur Ayam, Telur Bebek — Citra (Threads), Desy (Threads), Apsari (Threads), +2 lainnya [Bank BPD]',
    'posted',
    449170.00,   -- total debit (bank + cogs)
    449170.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  264500.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   184670.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 264500.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 184670.00);

  -- ── JE-2509-0006: 2025-09-22 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0006',
    '2025-09-22',
    'Penjualan Telur Ayam, Telur Bebek — Shanti Bintaran, Mbak Nita, Bulik Lili, +3 lainnya [Bank BPD]',
    'posted',
    476080.00,   -- total debit (bank + cogs)
    476080.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  280000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   196080.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 280000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 196080.00);

  -- ── JE-2509-0007: 2025-09-23 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0007',
    '2025-09-23',
    'Penjualan Telur Ayam, Telur Bebek — Irma Tk.Elektronik [Bank BPD]',
    'posted',
    88915.00,   -- total debit (bank + cogs)
    88915.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  52500.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   36415.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 52500.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 36415.00);

  -- ── JE-2509-0008: 2025-09-24 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0008',
    '2025-09-24',
    'Penjualan Telur Bebek, Telur Asin, Telur Asin Mentah, Telur Ayam — Wirosaban, Mirul - DISPETARU, Apsari (Threads), +5 lainnya [Bank BPD]',
    'posted',
    659110.00,   -- total debit (bank + cogs)
    659110.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  370000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   289110.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 370000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 289110.00);

  -- ── JE-2509-0009: 2025-09-25 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0009',
    '2025-09-25',
    'Penjualan Telur Asin, Telur Ayam — Mbak Putri, Cyntia [Bank BPD]',
    'posted',
    304830.00,   -- total debit (bank + cogs)
    304830.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  175000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   129830.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 175000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 129830.00);

  -- ── JE-2509-0010: 2025-09-26 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0010',
    '2025-09-26',
    'Penjualan Telur Ayam — Mbak Anti [Bank BPD]',
    'posted',
    62295.00,   -- total debit (bank + cogs)
    62295.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  37500.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   24795.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 37500.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 24795.00);

  -- ── JE-2509-0011: 2025-09-30 [BPD] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0011',
    '2025-09-30',
    'Penjualan Telur Ayam, Telur Bebek — Khalista Wirobrajan, Fita UPY [Bank BPD]',
    'posted',
    94530.00,   -- total debit (bank + cogs)
    94530.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bank_id,   'Penerimaan penjualan — Bank BPD',  55000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   39530.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 55000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 39530.00);

  -- ── JE-2509-0012: 2025-09-30 [CASH] ──────────────────────────
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by)
  VALUES (
    v_je_id,
    'JE-2509-0012',
    '2025-09-30',
    'Penjualan Telur Bebek — Bu Rini -DISPETARU [Cash]',
    'posted',
    79500.00,   -- total debit (bank + cogs)
    79500.00,   -- total credit (sales + inventory)
    v_admin_id
  );

  -- Lines:
  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash_id,   'Penerimaan penjualan — Cash',  45000.00, 0),
    (gen_random_uuid(), v_je_id, v_cogs,           'HPP penjualan telur',                   34500.00, 0),
    (gen_random_uuid(), v_je_id, v_sales,          'Pendapatan penjualan telur',             0, 45000.00),
    (gen_random_uuid(), v_je_id, v_inventory,      'Pengurangan stok telur',                 0, 34500.00);

  RAISE NOTICE 'Selesai: % journal entries dibuat.', 12;
END $$;

-- Verifikasi
SELECT journal_number, date, description, total_debit, total_credit, status
FROM journal_entries
WHERE journal_number LIKE 'JE-2509-%'
ORDER BY date, journal_number;

-- ================================================================
-- SEED: Journal Entries Penjualan Oktober 2025 – Februari 2026
-- Strategi: 1 entry per hari per metode pembayaran
-- Transaksi BELUM BAYAR / Proses: SKIP
--
-- coa yang dibutuhkan (pastikan ada di DB kamu):
--   10101 Cash
--   10102 Bank BPD
--   10103 Bank BPD       ← BARU (insert dulu jika belum ada)
--   10104 Bank BPD   ← BARU (insert dulu jika belum ada)
--   10105 Inventory - Eggs (sebelumnya 10104 di Sep, sesuaikan)
--   40101 Egg Sales
--   50101 Egg Purchase Cost
-- ================================================================

-- ── Tambah coa Bank baru jika belum ada ─────────────────────────
-- (Jalankan bagian ini manual dulu, atau skip jika sudah ada)
/*
INSERT INTO coa (id, subgroup_id, code, name, status)
SELECT gen_random_uuid(),
       (SELECT id FROM coa_subgroups WHERE name ILIKE '%current asset%' LIMIT 1),
       '10103', 'Bank BPD', 1
WHERE NOT EXISTS (SELECT 1 FROM coa WHERE code = '10103');

INSERT INTO coa (id, subgroup_id, code, name, status)
SELECT gen_random_uuid(),
       (SELECT id FROM coa_subgroups WHERE name ILIKE '%current asset%' LIMIT 1),
       '10104', 'Bank BPD', 1
WHERE NOT EXISTS (SELECT 1 FROM coa WHERE code = '10104');
*/

DO $$
DECLARE
  v_cash     UUID := (SELECT id FROM coa WHERE code = '10101');
  v_bpd      UUID := (SELECT id FROM coa WHERE code = '10102');
  v_inv      UUID := (SELECT id FROM coa WHERE code = '10104');  -- Inventory - Eggs
  v_sales    UUID := (SELECT id FROM coa WHERE code = '40101');  -- Egg Sales
  v_cogs     UUID := (SELECT id FROM coa WHERE code = '50101');  -- Egg Purchase Cost
  v_ar       UUID := (SELECT id FROM coa WHERE code = '10103');  -- Accounts Receivable (fallback)
  v_admin_id UUID := (SELECT id FROM users WHERE email = 'admin@example.com' LIMIT 1);
  v_je_id    UUID;
BEGIN

  -- ════════════════════════════════════════════════════════
  -- OKTOBER
  -- ════════════════════════════════════════════════════════

  -- ── JE-2510-0015: 2025-10-10 [BPD] (31 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0015', '2025-10-10',
    'Penjualan Telur Bebek, Telur Ayam, Telur Asin, Telur Asin Mentah — Citra (Threads) smti RO, Lala (Threads) Tugu, Fia - Dispetaru RO, +21 lainnya [Bank BPD]',
    'posted', 1793650, 1793650, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   1039000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               754650, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 1039000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 754650);

  -- ── JE-2510-0016: 2025-10-10 [CASH] (11 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0016', '2025-10-10',
    'Penjualan Telur Ayam, Telur Asin, Telur Bebek — ida - RO, Jaka - Kosan, Mbak Aris, +5 lainnya [Cash]',
    'posted', 1019748, 1019748, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   547899, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               471849, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 547899),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 471849);

  -- ── JE-2510-0017: 2025-10-13 [BPD] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0017', '2025-10-13',
    'Penjualan Telur Ayam, Telur Bebek — Threads Tamsis , Threads Godean, Khalista , +1 lainnya [Bank BPD]',
    'posted', 251500, 251500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   145000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               106500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 145000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 106500);

  -- ── JE-2510-0018: 2025-10-14 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0018', '2025-10-14',
    'Penjualan Telur Ayam, Telur Bebek — Ajeng, Erika, Mirul [Bank BPD]',
    'posted', 160500, 160500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   92500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               68000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 92500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 68000);

  -- ── JE-2510-0019: 2025-10-15 [BPD] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0019', '2025-10-15',
    'Penjualan Telur Ayam — Threads COD Tamsis, Umi [Bank BPD]',
    'posted', 172000, 172000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   100000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               72000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 100000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 72000);

  -- ── JE-2510-0020: 2025-10-16 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0020', '2025-10-16',
    'Penjualan Telur Bebek, Telur Ayam — Lutfi Optik Manding, Ida - RO, Shanti Rumah [Bank BPD]',
    'posted', 245375, 245375, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   142500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               102875, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 142500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 102875);

  -- ── JE-2510-0021: 2025-10-17 [BPD] (5 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0021', '2025-10-17',
    'Penjualan Telur Asin Mentah, Telur Asin, Telur Ayam, Telur Bebek — Ajeng, Bulik Lili , Mbak Nila, +2 lainnya [Bank BPD]',
    'posted', 259625, 259625, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   146500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               113125, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 146500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 113125);

  -- ── JE-2510-0022: 2025-10-20 [BPD] (12 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0022', '2025-10-20',
    'Penjualan Telur Ayam, Telur Asin, Telur Bebek — Mbak Emil Godean, Evi Ngasem - Threads, Mbak Aris, +6 lainnya [Bank BPD]',
    'posted', 819975, 819975, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   472500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               347475, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 472500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 347475);

  -- ── JE-2510-0023: 2025-10-21 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0023', '2025-10-21',
    'Penjualan Telur Ayam — Lala Dispetaru, Bu Yuli Dispetaru, Bu Rini Dispetaru [Bank BPD]',
    'posted', 127950, 127950, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   75000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               52950, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 75000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 52950);

  -- ── JE-2510-0024: 2025-10-21 [CASH] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0024', '2025-10-21',
    'Penjualan Telur Ayam — Ninda Rumah [Cash]',
    'posted', 42650, 42650, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   25000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               17650, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 25000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 17650);

  -- ── JE-2510-0025: 2025-10-23 [BPD] (9 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0025', '2025-10-23',
    'Penjualan Telur Ayam, Telur Asin, Telur Bebek, Telur Asin Mentah — Intan Dispetaru, mbak desy, Lola Rumah, +2 lainnya [Bank BPD]',
    'posted', 339370, 339370, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   191000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               148370, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 191000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 148370);

  -- ── JE-2510-0026: 2025-10-24 [BPD] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0026', '2025-10-24',
    'Penjualan Telur Bebek, Telur Ayam — Lala (Threads) Tugu [Bank BPD]',
    'posted', 95650, 95650, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   55000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               40650, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 55000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 40650);

  -- ── JE-2510-0027: 2025-10-25 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0027', '2025-10-25',
    'Penjualan Telur Ayam — Mas Ishom (450) [Bank BPD]',
    'posted', 1919250, 1919250, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   1125000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               794250, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 1125000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 794250);

  -- ── JE-2510-0028: 2025-10-25 [CASH] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0028', '2025-10-25',
    'Penjualan Telur Ayam, Telur Bebek — Mas Ishom (Bonus) [Cash]',
    'posted', 24390, 24390, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   12195, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               12195, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 12195),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 12195);

  -- ── JE-2510-0029: 2025-10-27 [BPD] (9 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0029', '2025-10-27',
    'Penjualan Telur Ayam, Ayam Kampung, Telur Bebek — Ida - RO, Irma  Tk.Elektronik - Threads, Mbak Emil Godean, +3 lainnya [Bank BPD]',
    'posted', 658000, 658000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   380000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               278000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 380000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 278000);

  -- ── JE-2510-0030: 2025-10-28 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0030', '2025-10-28',
    'Penjualan Telur Ayam — Endah Apotek Threads [Bank BPD]',
    'posted', 42500, 42500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   25000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               17500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 25000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 17500);

  -- ── JE-2510-0031: 2025-10-29 [BPD] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0031', '2025-10-29',
    'Penjualan Telur Ayam, Telur Asin — Desy Sardjito (Threads), Mbak Nila, Notaris Pakualaman [Bank BPD]',
    'posted', 327000, 327000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   190000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               137000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 190000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 137000);

  -- ── JE-2510-0032: 2025-10-30 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0032', '2025-10-30',
    'Penjualan Telur Ayam — Mirul Dispetaru, Lala Dispetaru, Toko Sayur Congcat Threads [Bank BPD]',
    'posted', 398970, 398970, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   225000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               173970, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 225000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 173970);

  -- ── JE-2510-0033: 2025-10-30 [CASH] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2510-0033', '2025-10-30',
    'Penjualan Telur Bebek — Aurora Jokteng Threads [Cash]',
    'posted', 79500, 79500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   45000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               34500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 45000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 34500);

  -- ════════════════════════════════════════════════════════
  -- NOVEMBER
  -- ════════════════════════════════════════════════════════

  -- ── JE-2511-0034: 2025-11-11 [BPD] (15 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0034', '2025-11-11',
    'Penjualan Ayam Kampung, Telur Asin, Telur Ayam, Telur Bebek — Shanti Rumah, Mirul Dispetaru, Astri , +9 lainnya [Bank BPD]',
    'posted', 1000375, 1000375, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   567500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               432875, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 567500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 432875);

  -- ── JE-2511-0035: 2025-11-11 [CASH] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0035', '2025-11-11',
    'Penjualan Telur Bebek, Telur Ayam — Rini Magang, Bu Rini [Cash]',
    'posted', 212500, 212500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   125000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               87500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 125000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 87500);

  -- ── JE-2511-0036: 2025-11-13 [BPD] (6 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0036', '2025-11-13',
    'Penjualan Telur Asin, Telur Ayam, Telur Bebek — Mbak ayuk jmc, Evi Ngasem, Khalita [Bank BPD]',
    'posted', 299628, 299628, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   169250, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               130378, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 169250),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 130378);

  -- ── JE-2511-0037: 2025-11-14 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0037', '2025-11-14',
    'Penjualan Telur Ayam, Telur Bebek — Yuliana Dispetaru, Siwi pajangan [Bank BPD]',
    'posted', 159700, 159700, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   92500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               67200, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 92500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 67200);

  -- ── JE-2511-0038: 2025-11-14 [CASH] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0038', '2025-11-14',
    'Penjualan Telur Ayam — Erika Dispetaru [Cash]',
    'posted', 64020, 64020, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   37500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               26520, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 37500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 26520);

  -- ── JE-2511-0039: 2025-11-17 [BPD] (5 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0039', '2025-11-17',
    'Penjualan Telur Ayam, Ayam Kampung — Mbak Emil Godean, Rizka Dispetaru, Gita SMP, +1 lainnya [Bank BPD]',
    'posted', 479000, 479000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   280000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               199000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 280000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 199000);

  -- ── JE-2511-0040: 2025-11-18 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0040', '2025-11-18',
    'Penjualan Telur Ayam — Mba Anti [Bank BPD]',
    'posted', 42500, 42500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   25000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               17500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 25000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 17500);

  -- ── JE-2511-0041: 2025-11-19 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0041', '2025-11-19',
    'Penjualan Telur Ayam — Owner [Bank BPD]',
    'posted', 17500, 17500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   8750, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               8750, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 8750),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 8750);

  -- ── JE-2511-0042: 2025-11-20 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0042', '2025-11-20',
    'Penjualan Telur Ayam — Khalista [Bank BPD]',
    'posted', 42500, 42500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   25000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               17500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 25000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 17500);

  -- ── JE-2511-0043: 2025-11-21 [BPD] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0043', '2025-11-21',
    'Penjualan Telur Ayam, Telur Bebek, Ayam Kampung — Andi Threads, Khalista, Mba Shanti , +1 lainnya [Bank BPD]',
    'posted', 320000, 320000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   185000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               135000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 185000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 135000);

  -- ── JE-2511-0044: 2025-11-25 [BPD] (10 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2511-0044', '2025-11-25',
    'Penjualan Telur Bebek, Telur Ayam, Ayam Kampung — Sidiq Threads, Mba Septi DPTR, Mbak Nia, +3 lainnya [Bank BPD]',
    'posted', 330550, 330550, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   186850, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               143700, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 186850),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 143700);

  -- ════════════════════════════════════════════════════════
  -- DESEMBER
  -- ════════════════════════════════════════════════════════

  -- ── JE-2512-0045: 2025-12-12 [BPD] (22 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0045', '2025-12-12',
    'Penjualan Telur Ayam, Telur Bebek, Ayam Kampung — Mba Astri DISPETRU, Siwi Gesikan , Mbak erika, +11 lainnya [Bank BPD]',
    'posted', 1733910, 1733910, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   985000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               748910, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 985000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 748910);

  -- ── JE-2512-0046: 2025-12-12 [CASH] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0046', '2025-12-12',
    'Penjualan Telur Ayam, Ayam Kampung — Mbak Nila, Mbak Lala, Bu Rini, +1 lainnya [Cash]',
    'posted', 221800, 221800, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   124000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               97800, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 124000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 97800);

  -- ── JE-2512-0047: 2025-12-14 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0047', '2025-12-14',
    'Penjualan Telur Ayam, Telur Asin — Threads Jakal, Mbak Gesti Dispetaru [Bank BPD]',
    'posted', 345610, 345610, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   195000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               150610, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 195000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 150610);

  -- ── JE-2512-0048: 2025-12-17 [BCA] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0048', '2025-12-17',
    'Penjualan Telur Bebek, Telur Ayam — Fia Dispetaru, Mbak Sari [Bank BPD]',
    'posted', 167800, 167800, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   95000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               72800, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 95000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 72800);

  -- ── JE-2512-0049: 2025-12-17 [MANDIRI] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0049', '2025-12-17',
    'Penjualan Telur Bebek, Telur Ayam, Ayam Kampung — Toefliana Mirota threads, Mas agustinus [Bank BPD]',
    'posted', 216955, 216955, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   122500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               94455, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 122500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 94455);

  -- ── JE-2512-0050: 2025-12-18 [BCA] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0050', '2025-12-18',
    'Penjualan Telur Ayam, Ayam Kampung — Ida [Bank BPD]',
    'posted', 147225, 147225, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   82500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               64725, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 82500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 64725);

  -- ── JE-2512-0051: 2025-12-18 [MANDIRI] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0051', '2025-12-18',
    'Penjualan Ayam Kampung — Apip Threads [Bank BPD]',
    'posted', 108000, 108000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   60000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               48000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 60000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 48000);

  -- ── JE-2512-0052: 2025-12-19 [BCA] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0052', '2025-12-19',
    'Penjualan Ayam Kampung, Telur Ayam — Sita [Bank BPD]',
    'posted', 49115, 49115, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   27500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               21615, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 27500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 21615);

  -- ── JE-2512-0053: 2025-12-22 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0053', '2025-12-22',
    'Penjualan Telur Ayam — Vika Atma [Bank BPD]',
    'posted', 132690, 132690, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   75000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               57690, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 75000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 57690);

  -- ── JE-2512-0054: 2025-12-24 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0054', '2025-12-24',
    'Penjualan Telur Asin — Mbak Gesti Dispetaru [Bank BPD]',
    'posted', 36000, 36000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   20000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               16000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 20000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 16000);

  -- ── JE-2512-0055: 2025-12-27 [BCA] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0055', '2025-12-27',
    'Penjualan Telur Asin, Telur Ayam — Mbak Rizka Dispetaru, Vika Atma [Bank BPD]',
    'posted', 740000, 740000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   411000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               329000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 411000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 329000);

  -- ── JE-2512-0056: 2025-12-31 [BCA] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0056', '2025-12-31',
    'Penjualan Telur Ayam, Ayam Kampung — Mba ayu dispetaru, Ida, Laksita [Bank BPD]',
    'posted', 242500, 242500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   137500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               105000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 137500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 105000);

  -- ── JE-2512-0057: 2025-12-31 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2512-0057', '2025-12-31',
    'Penjualan Telur Ayam, Ayam Kampung — Mas Maruf dispetaru, Wirosaban [Bank BPD]',
    'posted', 157000, 157000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   87500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               69500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 87500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 69500);

  -- ════════════════════════════════════════════════════════
  -- JANUARI
  -- ════════════════════════════════════════════════════════

  -- ── JE-2601-0058: 2026-01-01 [BCA] (10 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0058', '2026-01-01',
    'Penjualan Telur Bebek, Ayam Kampung, Telur Asin, Telur Ayam — Sanden Threads, Vika, Antik, +3 lainnya [Bank BPD]',
    'posted', 836910, 836910, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   477500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               359410, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 477500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 359410);

  -- ── JE-2601-0059: 2026-01-01 [BPD] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0059', '2026-01-01',
    'Penjualan Telur Ayam, Telur Asin — Mbak Aris, Mba Nila [Bank BPD]',
    'posted', 165250, 165250, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   95000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               70250, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 95000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 70250);

  -- ── JE-2601-0060: 2026-01-01 [CASH] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0060', '2026-01-01',
    'Penjualan Telur Ayam, Ayam Kampung — Brimob Threads [Cash]',
    'posted', 97000, 97000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   55000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               42000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 55000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 42000);

  -- ── JE-2601-0061: 2026-01-01 [MANDIRI] (5 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0061', '2026-01-01',
    'Penjualan Telur Ayam, Telur Asin, Ayam Kampung — Giwangan Spx, Mas Untung ppk, Toeflyana, +1 lainnya [Bank BPD]',
    'posted', 374865, 374865, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   212500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               162365, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 212500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 162365);

  -- ── JE-2601-0062: 2026-01-14 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0062', '2026-01-14',
    'Penjualan Telur Ayam — Vika [Bank BPD]',
    'posted', 123000, 123000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   75000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               48000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 75000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 48000);

  -- ── JE-2601-0063: 2026-01-14 [BPD] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0063', '2026-01-14',
    'Penjualan Telur Bebek, Telur Ayam — Intan [Bank BPD]',
    'posted', 47800, 47800, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   27000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               20800, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 27000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 20800);

  -- ── JE-2601-0064: 2026-01-14 [MANDIRI] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0064', '2026-01-14',
    'Penjualan Telur Bebek — Threads Hotel Tara [Bank BPD]',
    'posted', 110000, 110000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   60000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               50000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 60000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 50000);

  -- ── JE-2601-0065: 2026-01-16 [BCA] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0065', '2026-01-16',
    'Penjualan Telur Ayam — Yulia Perumahan, Threads Paskas [Bank BPD]',
    'posted', 246000, 246000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   150000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               96000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 150000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 96000);

  -- ── JE-2601-0066: 2026-01-19 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0066', '2026-01-19',
    'Penjualan Telur Asin — po Asin Ai Jawir  [Bank BPD]',
    'posted', 75000, 75000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   40000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               35000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 40000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 35000);

  -- ── JE-2601-0067: 2026-01-19 [MANDIRI] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0067', '2026-01-19',
    'Penjualan Telur Asin, Telur Ayam — po Asin Ayu Triana Alkid , po Asin Uthi Kotagede, po Asin Palagan, +1 lainnya [Bank BPD]',
    'posted', 286500, 286500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   157500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               129000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 157500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 129000);

  -- ── JE-2601-0068: 2026-01-20 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0068', '2026-01-20',
    'Penjualan Telur Ayam — Mbak Gesti [Bank BPD]',
    'posted', 82000, 82000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   50000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               32000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 50000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 32000);

  -- ── JE-2601-0069: 2026-01-21 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0069', '2026-01-21',
    'Penjualan Telur Ayam — Vika [Bank BPD]',
    'posted', 82000, 82000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   50000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               32000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 50000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 32000);

  -- ── JE-2601-0070: 2026-01-22 [BCA] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0070', '2026-01-22',
    'Penjualan Telur Ayam, Ayam Kampung, Telur Bebek — Ida, Hotel Tara [Bank BPD]',
    'posted', 505000, 505000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   285000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               220000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 285000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 220000);

  -- ── JE-2601-0071: 2026-01-26 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0071', '2026-01-26',
    'Penjualan Telur Ayam — mAS Agustinus [Bank BPD]',
    'posted', 62700, 62700, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   37500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               25200, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 37500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 25200);

  -- ── JE-2601-0072: 2026-01-27 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0072', '2026-01-27',
    'Penjualan Telur Ayam — Vika  [Bank BPD]',
    'posted', 126000, 126000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   75000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               51000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 75000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 51000);

  -- ── JE-2601-0073: 2026-01-27 [MANDIRI] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0073', '2026-01-27',
    'Penjualan Telur Ayam, Telur Bebek — Toeflyana, Toeflyana  [Bank BPD]',
    'posted', 145500, 145500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   82500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               63000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 82500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 63000);

  -- ── JE-2601-0074: 2026-01-28 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0074', '2026-01-28',
    'Penjualan Telur Ayam — Mbak gesti  [Bank BPD]',
    'posted', 63000, 63000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   37500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               25500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 37500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 25500);

  -- ── JE-2601-0075: 2026-01-28 [MANDIRI] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2601-0075', '2026-01-28',
    'Penjualan Telur Ayam — Mbak Antik [Bank BPD]',
    'posted', 42000, 42000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   25000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               17000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 25000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 17000);

  -- ════════════════════════════════════════════════════════
  -- FEBRUARI
  -- ════════════════════════════════════════════════════════

  -- ── JE-2602-0076: 2026-02-02 [BCA] (6 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0076', '2026-02-02',
    'Penjualan Telur Ayam, Ayam Kampung — Vika, Ida, Paseban Threads [Bank BPD]',
    'posted', 859800, 859800, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   510000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               349800, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 510000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 349800);

  -- ── JE-2602-0077: 2026-02-02 [BPD] (10 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0077', '2026-02-02',
    'Penjualan Telur Ayam, Telur Asin, Telur Asin Mentah — Mirul, Mbak Nila , Mbak Astri, +3 lainnya [Bank BPD]',
    'posted', 585170, 585170, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   332500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               252670, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 332500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 252670);

  -- ── JE-2602-0078: 2026-02-02 [CASH] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0078', '2026-02-02',
    'Penjualan Telur Ayam — Andi Tugu [Cash]',
    'posted', 167480, 167480, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   100000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               67480, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 100000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 67480);

  -- ── JE-2602-0079: 2026-02-02 [MANDIRI] (3 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0079', '2026-02-02',
    'Penjualan Telur Ayam, Telur Bebek — Toeflyana, Mbak ayu, Hotel Tara [Bank BPD]',
    'posted', 443620, 443620, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   250000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               193620, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 250000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 193620);

  -- ── JE-2602-0080: 2026-02-16 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0080', '2026-02-16',
    'Penjualan Telur Ayam — Paseban Threads [Bank BPD]',
    'posted', 175440, 175440, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   100000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               75440, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 100000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 75440);

  -- ── JE-2602-0081: 2026-02-16 [CASH] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0081', '2026-02-16',
    'Penjualan Telur Ayam — Swalayan Rizma [Cash]',
    'posted', 448100, 448100, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   250000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               198100, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 250000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 198100);

  -- ── JE-2602-0082: 2026-02-17 [BCA] (2 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0082', '2026-02-17',
    'Penjualan Telur Bebek, Telur Ayam — Srandakan threads  [Bank BPD]',
    'posted', 76930, 76930, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   42500, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               34430, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 42500),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 34430);

  -- ── JE-2602-0083: 2026-02-18 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0083', '2026-02-18',
    'Penjualan Telur Bebek — Aulia threads  [Bank BPD]',
    'posted', 82500, 82500, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   45000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               37500, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 45000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 37500);

  -- ── JE-2602-0084: 2026-02-20 [BCA] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0084', '2026-02-20',
    'Penjualan Ayam Kampung — Sita  [Bank BPD]',
    'posted', 162000, 162000, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   90000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               72000, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 90000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 72000);

  -- ── JE-2602-0085: 2026-02-20 [CASH] (4 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0085', '2026-02-20',
    'Penjualan Ayam Kampung, Telur Ayam — Mbak Ninda, Swalayan Rizma [Cash]',
    'posted', 860460, 860460, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_cash,   'Penerimaan — Cash',   485000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               375460, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 485000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 375460);

  -- ── JE-2602-0086: 2026-02-22 [BPD] (1 transaksi) ──
  v_je_id := gen_random_uuid();
  INSERT INTO journal_entries (id, journal_number, date, description, status, total_debit, total_credit, created_by) VALUES (
    v_je_id, 'JE-2602-0086', '2026-02-22',
    'Penjualan Telur Ayam — Mbak Nila  [Bank BPD]',
    'posted', 87720, 87720, v_admin_id);

  INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit) VALUES
    (gen_random_uuid(), v_je_id, v_bpd,   'Penerimaan — Bank BPD',   50000, 0),
    (gen_random_uuid(), v_je_id, v_cogs,      'HPP penjualan',               37720, 0),
    (gen_random_uuid(), v_je_id, v_sales,     'Pendapatan penjualan',         0, 50000),
    (gen_random_uuid(), v_je_id, v_inv,       'Pengurangan stok',             0, 37720);

  RAISE NOTICE 'Selesai: 72 journal entries Oktober–Februari dibuat.';
END $$;

-- ── Verifikasi ──────────────────────────────────────────────────
SELECT journal_number, date, LEFT(description,60) AS desc,
       total_debit, total_credit, status
FROM journal_entries
WHERE journal_number LIKE 'JE-25%' OR journal_number LIKE 'JE-26%'
ORDER BY date, journal_number;

-- ── Grand Total Check ───────────────────────────────────────────
SELECT
  SUM(CASE WHEN jl.coa_id = (SELECT id FROM coa WHERE code='40101') THEN jl.credit ELSE 0 END) AS total_revenue,
  SUM(CASE WHEN jl.coa_id = (SELECT id FROM coa WHERE code='50101') THEN jl.debit  ELSE 0 END) AS total_cogs,
  SUM(CASE WHEN jl.coa_id = (SELECT id FROM coa WHERE code='40101') THEN jl.credit ELSE 0 END)
  - SUM(CASE WHEN jl.coa_id = (SELECT id FROM coa WHERE code='50101') THEN jl.debit ELSE 0 END) AS profit
FROM journal_lines jl
INNER JOIN journal_entries je ON je.id = jl.journal_entry_id
WHERE je.status = 'posted'
  AND je.date >= '2025-10-01';