-- Migration down
-- [start][activities] =====================================================================================================
DROP INDEX IF EXISTS "activities_created_at_idx";
DROP INDEX IF EXISTS "activities_user_id_idx";
DROP INDEX IF EXISTS "activities_table_name_record_id_idx";
DROP TABLE IF EXISTS activities CASCADE;
-- [end][activities] =======================================================================================================

-- [start][rejected_images] ================================================================================================
DROP INDEX IF EXISTS "rejected_images_image_id_unique";
DROP TABLE IF EXISTS "rejected_images" CASCADE;
-- [end][rejected_images] ==================================================================================================

-- [start][accepted_images] ================================================================================================
DROP INDEX IF EXISTS "accepted_images_image_id_unique";
DROP TABLE IF EXISTS "accepted_images" CASCADE;
-- [end][accepted_images] ==================================================================================================

-- [start][images] =========================================================================================================
DROP INDEX IF EXISTS "images_deleted_at_idx";
DROP INDEX IF EXISTS "images_user_id_idx";
DROP TABLE IF EXISTS "images" CASCADE;
-- [end][images] ===========================================================================================================

-- [start][user_2fa_secrets] ===============================================================================================
DROP INDEX IF EXISTS "user_2fa_secrets_deleted_at_idx";
DROP INDEX IF EXISTS "user_2fa_secrets_user_id_idx";
DROP TABLE IF EXISTS user_2fa_secrets CASCADE;
-- [end][user_2fa_secrets] =================================================================================================

-- [start][sessions] =======================================================================================================
DROP INDEX IF EXISTS "sessions_deleted_at_idx";
DROP INDEX IF EXISTS "sessions_user_id_idx";
DROP TABLE IF EXISTS "sessions" CASCADE;
-- [end][sessions] =========================================================================================================

-- [start][bans] ===========================================================================================================
DROP INDEX IF EXISTS "bans_banned_by_idx";
DROP INDEX IF EXISTS "bans_user_id_idx";
DROP TABLE IF EXISTS bans CASCADE;
-- [end][bans] =============================================================================================================

-- [start][users] ============================================================================================================
DROP INDEX IF EXISTS "users_deleted_at_idx";
DROP INDEX IF EXISTS "users_email_unique";
DROP INDEX IF EXISTS "users_nickname_unique";
DROP TABLE IF EXISTS "users" CASCADE;
-- [end][users] ============================================================================================================

-- [start][reset_passwords] ==================================================================================================
DROP INDEX IF EXISTS "reset_passwords_email_unique";
DROP TABLE IF EXISTS reset_passwords CASCADE;
-- [end][reset_passwords] ====================================================================================================
