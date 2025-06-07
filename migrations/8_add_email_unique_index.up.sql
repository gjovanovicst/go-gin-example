-- Add unique index for email
ALTER TABLE `blog_auth` ADD UNIQUE KEY `email_unique` (`email`);