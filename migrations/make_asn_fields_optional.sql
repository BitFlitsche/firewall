-- Migration to make ASN fields optional
-- This migration makes RIR, Domain, and Country fields nullable in the ASN table

-- Make RIR field nullable
ALTER TABLE asns MODIFY COLUMN rir VARCHAR(20) NULL;

-- Make Domain field nullable  
ALTER TABLE asns MODIFY COLUMN domain VARCHAR(255) NULL;

-- Make Country field nullable
ALTER TABLE asns MODIFY COLUMN cc VARCHAR(2) NULL;

-- Add comment to document the change
ALTER TABLE asns COMMENT = 'ASN table with optional RIR, Domain, and Country fields'; 