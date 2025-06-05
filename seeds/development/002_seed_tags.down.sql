-- Rollback development seed data for tags table
DELETE FROM `blog_tag` WHERE `name` IN ('Technology', 'Programming', 'Go', 'Web Development', 'API', 'Testing');