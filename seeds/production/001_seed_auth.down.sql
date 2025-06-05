-- Rollback production seed data for auth table
DELETE FROM `blog_auth` WHERE `id` = 1;