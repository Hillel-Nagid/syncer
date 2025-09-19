-- Remove service records (handled here since they're inserted in this migration)
DELETE FROM services
WHERE name IN ('spotify', 'deezer');
-- Drop table
DROP TABLE IF EXISTS services;