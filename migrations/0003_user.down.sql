DROP INDEX CONCURRENTLY IF EXISTS "user_soft_delete_cnt";
DROP INDEX CONCURRENTLY IF EXISTS "user_phone_idx";
DROP INDEX CONCURRENTLY IF EXISTS "user_name_idx";
DROP INDEX CONCURRENTLY IF EXISTS "user_deleted_at_idx";
DROP TABLE IF EXISTS "user";