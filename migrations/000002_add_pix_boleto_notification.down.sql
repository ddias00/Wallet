-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_delivery_status;
DROP INDEX IF EXISTS idx_notifications_channel;
DROP INDEX IF EXISTS idx_notifications_event_type;
DROP INDEX IF EXISTS idx_notifications_user_id;

DROP INDEX IF EXISTS idx_boletos_created_at;
DROP INDEX IF EXISTS idx_boletos_due_date;
DROP INDEX IF EXISTS idx_boletos_status;
DROP INDEX IF EXISTS idx_boletos_code;
DROP INDEX IF EXISTS idx_boletos_transaction_id;
DROP INDEX IF EXISTS idx_boletos_user_id;

DROP INDEX IF EXISTS idx_pix_keys_created_at;
DROP INDEX IF EXISTS idx_pix_keys_status;
DROP INDEX IF EXISTS idx_pix_keys_type;
DROP INDEX IF EXISTS idx_pix_keys_key;
DROP INDEX IF EXISTS idx_pix_keys_user_id;

DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_type;

DROP INDEX IF EXISTS idx_payments_created_at;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_method;
DROP INDEX IF EXISTS idx_payments_transaction_id;
DROP INDEX IF EXISTS idx_payments_user_id;

-- Drop tables
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS boletos;
DROP TABLE IF EXISTS pix_keys;

-- Remove metadata column from transactions
ALTER TABLE transactions DROP COLUMN IF EXISTS metadata;
