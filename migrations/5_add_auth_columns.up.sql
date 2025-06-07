-- Add email and social login fields to auth table
ALTER TABLE `blog_auth` 
ADD COLUMN `email` varchar(255) DEFAULT '' COMMENT 'Email address',
ADD COLUMN `first_name` varchar(100) DEFAULT '' COMMENT 'First name',
ADD COLUMN `last_name` varchar(100) DEFAULT '' COMMENT 'Last name',
ADD COLUMN `avatar_url` varchar(500) DEFAULT '' COMMENT 'Avatar image URL',
ADD COLUMN `provider` varchar(50) DEFAULT 'local' COMMENT 'Auth provider (local, google, github, facebook)',
ADD COLUMN `provider_id` varchar(255) DEFAULT '' COMMENT 'Provider user ID',
ADD COLUMN `is_email_verified` boolean DEFAULT false COMMENT 'Email verification status',
ADD COLUMN `last_login` timestamp NULL COMMENT 'Last login timestamp',
ADD COLUMN `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT 'Account creation timestamp',
ADD COLUMN `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update timestamp';