-- Align legacy installs: domain.Car + AutoMigrate use owner_id; older SQL used client_id.
-- Safe no-op if owner_id already exists or client_id is absent.
BEGIN;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'cars' AND column_name = 'client_id'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'cars' AND column_name = 'owner_id'
  ) THEN
    ALTER TABLE cars RENAME COLUMN client_id TO owner_id;
  END IF;
END $$;

COMMIT;
