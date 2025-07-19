-- MySQL migration to fix is_cidr column name
-- Step 1: Add the correct column if it doesn't exist
ALTER TABLE ips ADD COLUMN IF NOT EXISTS is_c_id_r BOOLEAN NOT NULL DEFAULT FALSE;

-- Step 2: Copy data from the incorrect column to the correct column if needed
UPDATE ips SET is_c_id_r = is_cidr WHERE is_cidr IS NOT NULL;

-- Step 3: Drop the incorrect column
ALTER TABLE ips DROP COLUMN IF EXISTS is_cidr; 