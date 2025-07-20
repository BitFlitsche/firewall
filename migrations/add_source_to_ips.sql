-- Add source column to IPs table
ALTER TABLE ips ADD COLUMN source VARCHAR(50) NULL; 