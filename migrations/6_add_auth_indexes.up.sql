-- Update empty emails to unique values before adding unique constraint
UPDATE `blog_auth` SET `email` = CONCAT('user_', `id`, '@local.example') WHERE `email` = '' OR `email` IS NULL;