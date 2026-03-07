-- ============================================================
-- JOURNAL DUMMY DATA — Jan–Mar 2026
-- Ganti COMPANY_ID, FISCAL_PERIOD IDs, dan CREATED_BY
-- sebelum dijalankan!
-- ============================================================

BEGIN;

-- ────────────────────────────────────────────────────────
-- JANUARI 2026
-- ────────────────────────────────────────────────────────

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('e1ea73b5-fedb-4ef6-8a27-dbf0010500e4', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0001', 'general', '2026-01-01', 'Setoran modal awal pemilik', 'posted', 500000000, 500000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('0444fb7e-35ce-4778-84e1-ce09b1f0868a', 'e1ea73b5-fedb-4ef6-8a27-dbf0010500e4', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 500000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('81a472be-5a99-4dac-9648-79070b3e11e8', 'e1ea73b5-fedb-4ef6-8a27-dbf0010500e4', '30434024-e2ab-4986-8f92-03c4f22feb3e', 0, 500000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('6a694402-2de0-46fd-af3e-c5208e4a6175', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0002', 'general', '2026-01-02', 'Penerimaan pinjaman bank jangka panjang', 'posted', 200000000, 200000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('a10e500d-995b-4775-83bb-faa9010a261c', '6a694402-2de0-46fd-af3e-c5208e4a6175', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 200000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('9a4c0df0-5d7f-4840-b17a-ef1d22bb6ea8', '6a694402-2de0-46fd-af3e-c5208e4a6175', '4ab4c5f2-d2a6-4a6c-a251-d3bbe40efac8', 0, 200000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('91f4dc5f-0625-4830-ada1-6617403db595', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0003', 'general', '2026-01-03', 'Pembelian peralatan kantor', 'posted', 80000000, 80000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('3f2794a2-d8ff-4a8a-8515-e6d8caf9b298', '91f4dc5f-0625-4830-ada1-6617403db595', '12c94c1d-ed63-4387-948b-590c6d703719', 80000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('cce5dae2-4b1a-4f9e-9c6a-0964c147e68f', '91f4dc5f-0625-4830-ada1-6617403db595', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 80000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('44eaa603-c0d3-4524-a822-8ce5b44b9222', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0004', 'general', '2026-01-04', 'Pembelian kendaraan operasional', 'posted', 120000000, 120000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('3c745757-ef4c-4963-8e59-64e8ee219e00', '44eaa603-c0d3-4524-a822-8ce5b44b9222', '578541da-339e-4d84-9213-78dd7b63fcef', 120000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('c3f54fac-88cd-4eb8-81c2-9a463fb65885', '44eaa603-c0d3-4524-a822-8ce5b44b9222', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 120000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('6cf2b715-e46f-4b7f-a8bf-0a5edcc56ab1', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0005', 'general', '2026-01-05', 'Pembelian lisensi software akuntansi', 'posted', 15000000, 15000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('3dab43d9-66a2-4376-b2db-2c29b0db2a6b', '6cf2b715-e46f-4b7f-a8bf-0a5edcc56ab1', '6106c007-0586-4f95-9ea8-2b9386b9a467', 15000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ec98fdf4-8a36-4def-abdb-87635683e928', '6cf2b715-e46f-4b7f-a8bf-0a5edcc56ab1', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 15000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('5aa46064-ef20-4f17-aed2-2cac7516fc87', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0006', 'general', '2026-01-05', 'Pembayaran sewa gedung Jan-Mar 2026', 'posted', 18000000, 18000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('903a079d-7b06-4fb5-87aa-fd338c07510f', '5aa46064-ef20-4f17-aed2-2cac7516fc87', '9b78803f-69b0-4d88-924c-f0f8f58593eb', 18000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('bec3472c-e7e0-4a7f-845e-fccaa76597c3', '5aa46064-ef20-4f17-aed2-2cac7516fc87', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 18000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('cdf16564-0b8d-4aed-b3f6-562ee731d5da', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0007', 'general', '2026-01-08', 'Pembelian persediaan barang dagang dari supplier', 'posted', 150000000, 150000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('93ae3739-efcc-40aa-b386-5880a2175332', 'cdf16564-0b8d-4aed-b3f6-562ee731d5da', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 150000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('6e525f46-4fa2-4985-ad32-6610dea09f8c', 'cdf16564-0b8d-4aed-b3f6-562ee731d5da', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 0, 150000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('6165c700-b969-4113-95cb-f1a17b0064c7', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'RV-2026-0001', 'revenue', '2026-01-12', 'Penjualan barang dagang kepada pelanggan - tunai', 'posted', 85000000, 85000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('22d68ecf-f9a5-4654-820e-1326f1ce1134', '6165c700-b969-4113-95cb-f1a17b0064c7', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 85000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('af3acfda-b3a5-43c3-8244-57bf0feeb19d', '6165c700-b969-4113-95cb-f1a17b0064c7', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 85000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('29686e58-7824-4dcb-a307-95880324f648', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0008', 'general', '2026-01-12', 'HPP atas penjualan barang dagang 12 Jan', 'posted', 55000000, 55000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('823b9f9c-561f-445c-92c6-a39a2f71a484', '29686e58-7824-4dcb-a307-95880324f648', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 55000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('0bf5ef90-cad1-4e80-998e-ff14255b7f96', '29686e58-7824-4dcb-a307-95880324f648', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 55000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('1637ee32-41ac-4772-9977-ff0f0391a04d', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'RV-2026-0002', 'revenue', '2026-01-15', 'Penjualan barang kepada PT Maju Jaya - kredit 30 hari', 'posted', 62000000, 62000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('b8ae1743-266f-4304-a6e6-3bff74700d80', '1637ee32-41ac-4772-9977-ff0f0391a04d', '9728ea3d-9692-4f3c-946b-0322a0c4df7b', 62000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('78fa37f6-5b23-4856-9aef-930398301a71', '1637ee32-41ac-4772-9977-ff0f0391a04d', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 62000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('ca3e361f-b265-4570-9938-6a5e120706f6', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0009', 'general', '2026-01-15', 'HPP atas penjualan kredit 15 Jan', 'posted', 40000000, 40000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2d6eabcc-44cd-4553-80e3-0766ef3f93fc', 'ca3e361f-b265-4570-9938-6a5e120706f6', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 40000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('50cf3c5b-1cc0-4db7-bf93-90bb2f0f4225', 'ca3e361f-b265-4570-9938-6a5e120706f6', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 40000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('bb88aa5c-637b-4724-908a-d6ac6534a95d', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'RV-2026-0003', 'revenue', '2026-01-18', 'Pendapatan jasa konsultasi - PT Sejahtera', 'posted', 25000000, 25000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('933615e0-b656-4be0-9cb4-9b4b08f4507e', 'bb88aa5c-637b-4724-908a-d6ac6534a95d', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 25000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('41758128-0187-4f3f-a9de-6b0d25c7910c', 'bb88aa5c-637b-4724-908a-d6ac6534a95d', '4c548201-4c95-4fca-a056-63808b2813e5', 0, 25000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('8a268630-8363-42ad-be90-f9f82fe81a1f', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'EX-2026-0001', 'expense', '2026-01-31', 'Pembayaran gaji karyawan Januari 2026', 'posted', 35000000, 35000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('973f2013-d420-4ed9-8b73-2ce960cb78cb', '8a268630-8363-42ad-be90-f9f82fe81a1f', 'ff64e581-e5a2-406c-9647-4f7d08c1d447', 35000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('3308c803-a61f-4377-b774-48edfefbbd3c', '8a268630-8363-42ad-be90-f9f82fe81a1f', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 35000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('e970a2af-15e1-4ed9-8da1-78a874acf521', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'EX-2026-0002', 'expense', '2026-01-31', 'Pembayaran listrik & air Januari 2026', 'posted', 3500000, 3500000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('a1ad68e9-1c4a-4d51-a924-36517fbcef7e', 'e970a2af-15e1-4ed9-8da1-78a874acf521', '721086f5-8903-4cd6-bbb0-85018755643c', 3500000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('f3cf1851-92a5-4e9e-af4d-9922ce32e2bc', 'e970a2af-15e1-4ed9-8da1-78a874acf521', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 3500000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('e12e03f0-0fcd-44b1-a695-5c3c14300cfd', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'EX-2026-0003', 'expense', '2026-01-31', 'Pengakuan beban sewa Januari 2026', 'posted', 6000000, 6000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('fce555c0-d3b0-4fb9-b7b0-4aada95d8a63', 'e12e03f0-0fcd-44b1-a695-5c3c14300cfd', '285a8fba-2941-4848-8cb6-ce2e245609b2', 6000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('c3d766ae-2718-4ef3-a119-32f6e089d97d', 'e12e03f0-0fcd-44b1-a695-5c3c14300cfd', '9b78803f-69b0-4d88-924c-f0f8f58593eb', 0, 6000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('762eaf4f-5e4b-4c7a-98fd-1c7f948a194e', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'JE-2026-0010', 'general', '2026-01-31', 'Pembayaran sebagian hutang ke supplier', 'posted', 80000000, 80000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ce7d1dfd-0d8e-4ced-b4ee-d0f6b6eb894e', '762eaf4f-5e4b-4c7a-98fd-1c7f948a194e', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 80000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('fe51989e-34cc-4126-9618-821a934392e1', '762eaf4f-5e4b-4c7a-98fd-1c7f948a194e', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 80000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('a8e81014-8ecd-4813-b996-ffbab71d56b9', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'EX-2026-0004', 'expense', '2026-01-31', 'Biaya administrasi bank Januari 2026', 'posted', 250000, 250000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('190e38cf-a80d-429a-ab79-648c1543fa60', 'a8e81014-8ecd-4813-b996-ffbab71d56b9', 'eb4ae642-6df4-42bc-9e51-e4b0ba5d2df1', 250000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('f21fd09b-b9c3-4797-970b-f878ccf7c3c4', 'a8e81014-8ecd-4813-b996-ffbab71d56b9', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 250000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('a02b60c6-a9bd-4c60-94fb-dc37a5223959', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '48b48e6e-e36c-4bf3-b91f-d56e505d59e7', 'EX-2026-0005', 'expense', '2026-01-20', 'Pembelian perlengkapan kantor', 'posted', 2500000, 2500000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('b8c407d6-654d-4c0b-9b78-ac64b9a3330f', 'a02b60c6-a9bd-4c60-94fb-dc37a5223959', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 2500000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('33bfe375-114d-43d5-b192-0681d4ec863e', 'a02b60c6-a9bd-4c60-94fb-dc37a5223959', 'c9be7f73-d70c-4f81-be8f-9c721f653987', 0, 2500000);

-- ────────────────────────────────────────────────────────
-- FEBRUARI 2026
-- ────────────────────────────────────────────────────────

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('34c17a11-c084-4dd5-adcd-9bcde672a945', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0011', 'general', '2026-02-05', 'Penerimaan pelunasan piutang dari PT Maju Jaya', 'posted', 62000000, 62000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('d473b17d-1476-4758-bda6-ca3cd48de8e7', '34c17a11-c084-4dd5-adcd-9bcde672a945', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 62000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('278d8b86-7d8b-4092-b142-1e1835615c0d', '34c17a11-c084-4dd5-adcd-9bcde672a945', '9728ea3d-9692-4f3c-946b-0322a0c4df7b', 0, 62000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('c1df9fe9-e5d3-48f1-bfd3-14798bed5c03', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0012', 'general', '2026-02-06', 'Pembelian persediaan barang dagang', 'posted', 180000000, 180000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('46d9dd01-97a6-4ee0-a91c-c57c4b9fc4dd', 'c1df9fe9-e5d3-48f1-bfd3-14798bed5c03', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 180000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ec6b0f79-84f2-4a0a-b085-27cbfa3d3ae6', 'c1df9fe9-e5d3-48f1-bfd3-14798bed5c03', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 0, 180000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('651fe158-17d8-4ad6-ae61-429578f8f5ff', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'RV-2026-0004', 'revenue', '2026-02-10', 'Penjualan barang dagang - tunai', 'posted', 110000000, 110000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('0e9e5dcf-faa0-46b9-90cc-dc0d0d9ec93f', '651fe158-17d8-4ad6-ae61-429578f8f5ff', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 110000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('33723d17-8673-4273-b6d5-6cb572288dd0', '651fe158-17d8-4ad6-ae61-429578f8f5ff', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 110000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('d601aa56-39d5-4253-8c59-dcc95bfb82ce', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0013', 'general', '2026-02-10', 'HPP atas penjualan 10 Feb', 'posted', 71500000, 71500000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('e6220e3c-a7ac-492e-b11d-af6cbd891c41', 'd601aa56-39d5-4253-8c59-dcc95bfb82ce', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 71500000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('dad94d3b-d6ec-444d-81fd-689e94d22509', 'd601aa56-39d5-4253-8c59-dcc95bfb82ce', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 71500000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('242cfce2-546b-4c2e-9894-75fedd3803c8', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'RV-2026-0005', 'revenue', '2026-02-14', 'Penjualan kepada CV Berkah Abadi - kredit', 'posted', 78000000, 78000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('5ff1ecb1-c8c7-4383-9624-ed06704475fc', '242cfce2-546b-4c2e-9894-75fedd3803c8', '9728ea3d-9692-4f3c-946b-0322a0c4df7b', 78000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2448660f-77b9-4451-a972-9b409bcc4b3c', '242cfce2-546b-4c2e-9894-75fedd3803c8', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 78000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('510ff873-26be-4ca3-a944-dc7442d9d486', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0014', 'general', '2026-02-14', 'HPP atas penjualan kredit 14 Feb', 'posted', 50700000, 50700000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('e4b95276-35ee-456f-9bdf-17c9035dabd0', '510ff873-26be-4ca3-a944-dc7442d9d486', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 50700000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('acaa695f-7dc7-4320-8156-94bd9fd1dd34', '510ff873-26be-4ca3-a944-dc7442d9d486', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 50700000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('1a003033-e569-400d-ab62-731abe43c4ae', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'RV-2026-0006', 'revenue', '2026-02-16', 'Pendapatan jasa konsultasi - PT Maju Jaya', 'posted', 32000000, 32000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('52852ecb-dd41-41b8-8880-9df4bc097c0a', '1a003033-e569-400d-ab62-731abe43c4ae', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 32000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('d59edb2b-375d-46d6-9014-1b4c12604367', '1a003033-e569-400d-ab62-731abe43c4ae', '4c548201-4c95-4fca-a056-63808b2813e5', 0, 32000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('7ff351bb-f2ef-472c-967f-e03260b7760a', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0015', 'general', '2026-02-20', 'Pendapatan bunga deposito Februari 2026', 'posted', 1500000, 1500000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('1952bebe-a8aa-48e8-87cd-f4869dd5e4a8', '7ff351bb-f2ef-472c-967f-e03260b7760a', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 1500000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ee1e7955-6c39-4b23-9470-ca0398fc9571', '7ff351bb-f2ef-472c-967f-e03260b7760a', '0a948478-4b70-4597-8ee7-bb94226c058e', 0, 1500000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('d0381abe-cf9c-4e47-9bd5-6ceb1e8ebf53', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'EX-2026-0006', 'expense', '2026-02-28', 'Pembayaran gaji karyawan Februari 2026', 'posted', 37000000, 37000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('03460e82-b224-47b9-864d-4cfa85756e49', 'd0381abe-cf9c-4e47-9bd5-6ceb1e8ebf53', 'ff64e581-e5a2-406c-9647-4f7d08c1d447', 37000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('27c15915-66a2-40ef-91d5-3cc132956a94', 'd0381abe-cf9c-4e47-9bd5-6ceb1e8ebf53', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 37000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('5ac57d2e-be12-4475-a8d3-3c1f3921dd24', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'EX-2026-0007', 'expense', '2026-02-28', 'Pembayaran listrik & air Februari 2026', 'posted', 3800000, 3800000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('58bca563-85ab-4948-91c0-7afe7231cb9e', '5ac57d2e-be12-4475-a8d3-3c1f3921dd24', '721086f5-8903-4cd6-bbb0-85018755643c', 3800000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('a0de5721-996f-4155-986d-32d9eb9e1737', '5ac57d2e-be12-4475-a8d3-3c1f3921dd24', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 3800000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('48e1e7cf-344b-4a3d-9ba5-1ef3e0d16b99', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'EX-2026-0008', 'expense', '2026-02-28', 'Pengakuan beban sewa Februari 2026', 'posted', 6000000, 6000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('bd81f470-0ab7-4bc1-bf64-1693b2b9dbf2', '48e1e7cf-344b-4a3d-9ba5-1ef3e0d16b99', '285a8fba-2941-4848-8cb6-ce2e245609b2', 6000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('825a3ef9-91e3-4775-81c0-214f66612549', '48e1e7cf-344b-4a3d-9ba5-1ef3e0d16b99', '9b78803f-69b0-4d88-924c-f0f8f58593eb', 0, 6000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('141f6e49-5cb2-4088-b596-0f393dea263f', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0016', 'general', '2026-02-28', 'Pembayaran hutang ke supplier - Februari', 'posted', 100000000, 100000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('d9fbb75f-c59a-420b-a3ea-61c6d519e2e3', '141f6e49-5cb2-4088-b596-0f393dea263f', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 100000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ea49f6ff-9418-4968-8170-66f149ecb5ca', '141f6e49-5cb2-4088-b596-0f393dea263f', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 100000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('2e585ff0-f3e6-4776-a94d-5664dead8633', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'EX-2026-0009', 'expense', '2026-02-28', 'Biaya administrasi bank Februari 2026', 'posted', 275000, 275000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2a00eab7-9f14-48ea-ad58-bcc18ffad1f9', '2e585ff0-f3e6-4776-a94d-5664dead8633', 'eb4ae642-6df4-42bc-9e51-e4b0ba5d2df1', 275000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('778343d7-6b3c-4aad-a397-07d361384ad8', '2e585ff0-f3e6-4776-a94d-5664dead8633', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 275000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('f27addc5-44f6-4a83-b087-4c8d098b1328', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'EX-2026-0010', 'expense', '2026-02-15', 'Pembelian perlengkapan kantor Februari', 'posted', 1800000, 1800000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('cb5b623a-703b-4444-9564-6939c38ba2f7', 'f27addc5-44f6-4a83-b087-4c8d098b1328', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 1800000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('8b1e7459-7086-46e1-b926-6c297011b8b4', 'f27addc5-44f6-4a83-b087-4c8d098b1328', 'c9be7f73-d70c-4f81-be8f-9c721f653987', 0, 1800000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('ad6b3682-7383-447c-b096-31a1e097dcbd', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', '6ca66ef4-a9d8-417a-ad8b-bf1a02526847', 'JE-2026-0017', 'general', '2026-02-28', 'Penyusutan peralatan & kendaraan Februari 2026', 'posted', 3333333, 3333333, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('0ef58f6b-5283-408f-855e-b60ee7f38bed', 'ad6b3682-7383-447c-b096-31a1e097dcbd', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 3333333, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('40214888-e1fb-4da3-ae25-f2fd646067c8', 'ad6b3682-7383-447c-b096-31a1e097dcbd', 'f254867c-210f-4569-8a1c-87201e06a4e3', 0, 3333333);

-- ────────────────────────────────────────────────────────
-- MARET 2026
-- ────────────────────────────────────────────────────────

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('7d5b3731-65c9-42d8-be92-ad2aca2f487c', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0018', 'general', '2026-03-05', 'Penerimaan pelunasan piutang dari CV Berkah Abadi', 'posted', 78000000, 78000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('f76399be-de17-4513-92b9-a7e4ee8fd7fe', '7d5b3731-65c9-42d8-be92-ad2aca2f487c', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 78000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2c143e3e-1d8c-4438-bcff-195359c44348', '7d5b3731-65c9-42d8-be92-ad2aca2f487c', '9728ea3d-9692-4f3c-946b-0322a0c4df7b', 0, 78000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('cbcaa499-2f17-43cd-9a60-e13e1b849a67', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0019', 'general', '2026-03-07', 'Pembelian persediaan barang dagang Maret', 'posted', 200000000, 200000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2e463737-be49-4dbd-b702-09e53b0a31f7', 'cbcaa499-2f17-43cd-9a60-e13e1b849a67', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 200000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('2f9ee297-c3e0-461b-8d81-6de31d8a81e6', 'cbcaa499-2f17-43cd-9a60-e13e1b849a67', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 0, 200000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('ebc26c2b-1fee-4b90-9745-8109edd10564', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'RV-2026-0007', 'revenue', '2026-03-08', 'Penjualan barang dagang - tunai', 'posted', 95000000, 95000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('6d12c87c-cc28-4e3a-9f3c-56cbfafcf316', 'ebc26c2b-1fee-4b90-9745-8109edd10564', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 95000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('d5a66811-b06e-49df-9c13-b81802788779', 'ebc26c2b-1fee-4b90-9745-8109edd10564', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 95000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('4a7df870-59a3-40a2-a4a9-cb2cfe1ea293', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0020', 'general', '2026-03-08', 'HPP atas penjualan 8 Mar', 'posted', 61750000, 61750000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('6a16c768-c045-4f0f-88bc-5aeb523de1ac', '4a7df870-59a3-40a2-a4a9-cb2cfe1ea293', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 61750000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('16ed3c92-92d7-49ef-9b98-31f131ce7978', '4a7df870-59a3-40a2-a4a9-cb2cfe1ea293', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 61750000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('5cd05df0-e1b7-485c-ad51-1c022fc8f3ee', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'RV-2026-0008', 'revenue', '2026-03-12', 'Penjualan kepada PT Nusantara - kredit 30 hari', 'posted', 92000000, 92000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('b92072f5-ce14-48cb-a9fe-1e59991b97ab', '5cd05df0-e1b7-485c-ad51-1c022fc8f3ee', '9728ea3d-9692-4f3c-946b-0322a0c4df7b', 92000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('928b89c0-9e00-456a-b5a3-c82abee0b1ae', '5cd05df0-e1b7-485c-ad51-1c022fc8f3ee', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 92000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('7068b120-840d-4694-9830-cc9fb099af1d', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0021', 'general', '2026-03-12', 'HPP atas penjualan kredit 12 Mar', 'posted', 59800000, 59800000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('18fcb471-ecb3-4ce3-9159-96ac82269386', '7068b120-840d-4694-9830-cc9fb099af1d', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 59800000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('746ed7ce-23b5-4f22-aba0-6e55d7813f50', '7068b120-840d-4694-9830-cc9fb099af1d', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 59800000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('d82c3bd2-876e-4660-b366-6f965400c4bb', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'RV-2026-0009', 'revenue', '2026-03-20', 'Penjualan barang dagang - tunai (batch 2)', 'posted', 55000000, 55000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('6704616b-5ca4-4f4c-b7b2-2a0ecef99939', 'd82c3bd2-876e-4660-b366-6f965400c4bb', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 55000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('e5d4f87b-8af5-4fac-9226-187c0249e621', 'd82c3bd2-876e-4660-b366-6f965400c4bb', 'a11275ac-81df-4e8f-82aa-84476e963955', 0, 55000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('fd76d7c9-668d-446e-98a0-4a690ca4ede2', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0022', 'general', '2026-03-20', 'HPP atas penjualan 20 Mar', 'posted', 35750000, 35750000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('389522b6-e413-41e5-b777-59d66e686830', 'fd76d7c9-668d-446e-98a0-4a690ca4ede2', '0382ee9c-05fd-47c2-ac52-7002dd374f4f', 35750000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('94066c68-e6bf-4824-9623-cb64f904dbb0', 'fd76d7c9-668d-446e-98a0-4a690ca4ede2', 'fa44c46d-f3ba-4fbb-b6f9-0cc91214a57a', 0, 35750000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('0ed85099-ac75-4318-a30f-153c39da79b2', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'RV-2026-0010', 'revenue', '2026-03-18', 'Pendapatan jasa konsultasi Q1 - PT Sejahtera', 'posted', 45000000, 45000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('a03e9cbd-f7d6-4c26-9194-5ccc945e4aaa', '0ed85099-ac75-4318-a30f-153c39da79b2', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 45000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('205a7ebc-ec94-4b13-863b-e002402e995f', '0ed85099-ac75-4318-a30f-153c39da79b2', '4c548201-4c95-4fca-a056-63808b2813e5', 0, 45000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('bef8c87f-2fd1-423b-b934-5c37fbfc22a5', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0023', 'general', '2026-03-25', 'Pendapatan bunga tabungan Maret 2026', 'posted', 2000000, 2000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('4494ebc8-c2ff-4a8b-b794-706a39260ec7', 'bef8c87f-2fd1-423b-b934-5c37fbfc22a5', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 2000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('d0dfaf64-0021-4b06-b3d1-65aa5c8d4c3e', 'bef8c87f-2fd1-423b-b934-5c37fbfc22a5', '0a948478-4b70-4597-8ee7-bb94226c058e', 0, 2000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('2b14d316-a9a1-48ec-8bf4-929db7c8f256', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'EX-2026-0012', 'expense', '2026-03-31', 'Pembayaran gaji karyawan Maret 2026', 'posted', 38500000, 38500000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('3f799223-a2ab-4a52-a59e-315ae3bdcda0', '2b14d316-a9a1-48ec-8bf4-929db7c8f256', 'ff64e581-e5a2-406c-9647-4f7d08c1d447', 38500000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('df2e8315-6ba7-4b86-b067-351a6c97de05', '2b14d316-a9a1-48ec-8bf4-929db7c8f256', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 38500000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('445d03ad-1514-45f1-ae25-2c7a967df876', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'EX-2026-0013', 'expense', '2026-03-31', 'Pembayaran listrik & air Maret 2026', 'posted', 4200000, 4200000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('95911033-b571-4af2-af72-8e471f1a7a7b', '445d03ad-1514-45f1-ae25-2c7a967df876', '721086f5-8903-4cd6-bbb0-85018755643c', 4200000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('7be9b804-7ca0-403f-8c93-b7531eeaa1b1', '445d03ad-1514-45f1-ae25-2c7a967df876', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 4200000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('86597be2-df21-4e15-9acc-85a66996a391', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'EX-2026-0014', 'expense', '2026-03-31', 'Pengakuan beban sewa Maret 2026 (terakhir)', 'posted', 6000000, 6000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('bbfa55de-8d8b-4c30-9a27-be8e797bc6b0', '86597be2-df21-4e15-9acc-85a66996a391', '285a8fba-2941-4848-8cb6-ce2e245609b2', 6000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('7d78ed30-2ec0-4c43-92c8-73ce93b6339e', '86597be2-df21-4e15-9acc-85a66996a391', '9b78803f-69b0-4d88-924c-f0f8f58593eb', 0, 6000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('7efabdbe-f44b-45d1-9d82-e233d58af7e3', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0024', 'general', '2026-03-31', 'Pembayaran hutang ke supplier - Maret', 'posted', 130000000, 130000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('21f0fd0a-1745-490d-8464-90c4ec8309e9', '7efabdbe-f44b-45d1-9d82-e233d58af7e3', '201f8fb6-8a12-4e16-907f-746a1a3d8b0a', 130000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('ce9af93b-7f7a-4160-b0a7-ded422ea7665', '7efabdbe-f44b-45d1-9d82-e233d58af7e3', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 130000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('122be320-9436-4a35-b551-d5dfc297c843', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'EX-2026-0015', 'expense', '2026-03-31', 'Biaya administrasi bank Maret 2026', 'posted', 300000, 300000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('9678df4d-6c4d-442a-b0bc-78801e4d5923', '122be320-9436-4a35-b551-d5dfc297c843', 'eb4ae642-6df4-42bc-9e51-e4b0ba5d2df1', 300000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('30d5ddd9-c3c4-4a3c-9a31-04b8bebb19bd', '122be320-9436-4a35-b551-d5dfc297c843', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 300000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('479ae69f-4ba5-412c-a00b-429aa1c48ebc', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'EX-2026-0016', 'expense', '2026-03-22', 'Pembelian perlengkapan kantor Maret', 'posted', 2200000, 2200000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('91b7afa5-130f-4754-9665-1a2befe2df62', '479ae69f-4ba5-412c-a00b-429aa1c48ebc', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 2200000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('987b6dfb-c974-4035-b695-c35d15ffda9b', '479ae69f-4ba5-412c-a00b-429aa1c48ebc', 'c9be7f73-d70c-4f81-be8f-9c721f653987', 0, 2200000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('e1e75ad0-3be1-4222-bad2-60fd7183fa76', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0025', 'general', '2026-03-31', 'Penarikan dana oleh pemilik (drawings)', 'posted', 20000000, 20000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('845443f2-8e73-4e91-880c-1f33c8e0697c', 'e1e75ad0-3be1-4222-bad2-60fd7183fa76', '991ec2bc-5b08-4490-94d4-553ecaccda49', 20000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('601b687c-daea-41ed-a948-ca17b953da70', 'e1e75ad0-3be1-4222-bad2-60fd7183fa76', 'f241874f-6bd5-4c5f-9e79-017297b9dfe0', 0, 20000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('d0807357-6b77-4b29-bb1a-b09f61db0fcd', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0026', 'general', '2026-03-31', 'Pencadangan pajak penghasilan Q1 2026', 'posted', 12000000, 12000000, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('580ac28c-12a6-493f-a16b-9969fa59f0f7', 'd0807357-6b77-4b29-bb1a-b09f61db0fcd', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 12000000, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('275a1742-b3aa-49f6-9f4e-fd6687abe6fc', 'd0807357-6b77-4b29-bb1a-b09f61db0fcd', '46ca12cf-21b5-4f8a-8c1a-93ffc2b28860', 0, 12000000);

INSERT INTO public.journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_by)
VALUES ('60a4c897-9acc-4b73-9b1b-e95cd074929b', '20d5fffc-f812-4a6a-b920-0a0a0a1f3f35', 'c52e5ce2-3f5e-4278-8683-c77ff80dd274', 'JE-2026-0027', 'general', '2026-03-31', 'Penyusutan peralatan & kendaraan Maret 2026', 'posted', 3333333, 3333333, 'eaf906de-701e-4349-b002-69e7e6d3c677');
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('5ce56164-0ed5-4dac-942d-d8a702af178b', '60a4c897-9acc-4b73-9b1b-e95cd074929b', 'ed45e579-970f-4f40-b523-5ecf4cea7732', 3333333, 0);
INSERT INTO public.journal_lines (id, journal_entry_id, coa_id, debit, credit)
VALUES ('0d42a3a8-bc2f-4536-982b-375c4425bc32', '60a4c897-9acc-4b73-9b1b-e95cd074929b', 'f254867c-210f-4569-8a1c-87201e06a4e3', 0, 3333333);

COMMIT;