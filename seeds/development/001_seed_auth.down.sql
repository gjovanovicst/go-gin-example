-- Rollback development seed data for auth table
DELETE FROM `blog_auth` WHERE `id` IN (1, 2, 3);