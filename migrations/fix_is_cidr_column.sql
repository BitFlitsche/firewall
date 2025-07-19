-- Migration to fix is_cidr column name
-- This script will ensure only the correct column is used for CIDR flags

-- Step 1: Add the correct column if it doesn't exist
ALTER TABLE ips ADD COLUMN IF NOT EXISTS is_c_id_r BOOLEAN DEFAULT FALSE;

-- Step 2: Copy data from the incorrect column to the correct column if needed
UPDATE ips SET is_c_id_r = is_cidr WHERE is_cidr IS NOT NULL;

-- Step 3: Drop the incorrect column
ALTER TABLE ips DROP COLUMN IF EXISTS is_cidr;

-- Step 4: Ensure the correct column has the right default and constraints
ALTER TABLE ips ALTER COLUMN is_c_id_r SET DEFAULT FALSE;
ALTER TABLE ips ALTER COLUMN is_c_id_r SET NOT NULL; 