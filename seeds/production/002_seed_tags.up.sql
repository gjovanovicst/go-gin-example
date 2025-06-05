-- Production seed data for tags table - essential categories only
INSERT INTO `blog_tag` (`name`, `created_on`, `created_by`, `state`) VALUES 
('General', UNIX_TIMESTAMP(), 'system', 1);