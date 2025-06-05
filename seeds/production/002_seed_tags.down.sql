-- Rollback production seed data for tags table
DELETE FROM `blog_tag` WHERE `name` = 'General';