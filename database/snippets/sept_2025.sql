-- =============================================================
-- SEED: Journal Entries — Penjualan Telur September 2025
-- Strategi: 1 entry per hari per metode pembayaran
-- Setiap entry: 4 lines (Bank/Cash debit, HPP debit, Sales credit, Inventory credit)
--
-- PENTING: Jalankan query berikut dulu untuk dapat chart_of_accounts IDs kamu:
-- SELECT id, code, name FROM chart_of_accounts WHERE code IN
--   ('10101','10102','10104','40101','50101');
-- =============================================================

DO $$
DECLARE
  -- chart_of_accounts IDs — ambil dari database kamu
  v_cash_id    UUID := (SELECT id FROM chart_of_accounts WHERE code = '10101');  -- Cash
  v_bank_id    UUID := (SELECT id FROM chart_of_accounts WHERE code = '10102');  -- Bank
  v_inventory  UUID := (SELECT id FROM chart_of_accounts WHERE code = '10104');  -- Inventory - Eggs
  v_sales      UUID := (SELECT id FROM chart_of_accounts WHERE code = '40101');  -- Egg Sales
  v_cogs       UUID := (SELECT id FROM chart_of_accounts WHERE code = '50101');  -- Egg Purchase Cost
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