-- Drop trigger and function
DROP TRIGGER IF EXISTS update_user_services_updated_at ON user_services;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- Drop indexes
DROP INDEX IF EXISTS idx_user_services_next_sync;
DROP INDEX IF EXISTS idx_user_services_service_user;
DROP INDEX IF EXISTS idx_user_services_updated_at;
-- Drop table
DROP TABLE IF EXISTS user_services;