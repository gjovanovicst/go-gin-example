-- Remove social login fields from auth table
ALTER TABLE `blog_auth` 
DROP COLUMN `email`,
DROP COLUMN `first_name`,
DROP COLUMN `last_name`,
DROP COLUMN `avatar_url`,
DROP COLUMN `provider`,
DROP COLUMN `provider_id`,
DROP COLUMN `is_email_verified`,
DROP COLUMN `last_login`,
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`;