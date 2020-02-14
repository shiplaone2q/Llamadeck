-- Drop API key related indexes and tables first
DROP INDEX IF EXISTS idx_key_permissions_instance_id;
DROP INDEX IF EXISTS idx_api_keys_expires_at;
