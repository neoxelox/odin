DROP INDEX CONCURRENTLY IF EXISTS "otp_soft_delete_cnt";
DROP INDEX CONCURRENTLY IF EXISTS "otp_asset_idx";
DROP INDEX CONCURRENTLY IF EXISTS "otp_deleted_at_idx";
DROP TABLE IF EXISTS "otp";