-- Development seed data for tags table
INSERT INTO `blog_tag` (`name`, `created_on`, `created_by`, `state`) VALUES 
('Technology', UNIX_TIMESTAMP(), 'system', 1),
('Programming', UNIX_TIMESTAMP(), 'system', 1),
('Go', UNIX_TIMESTAMP(), 'system', 1),
('Web Development', UNIX_TIMESTAMP(), 'system', 1),
('API', UNIX_TIMESTAMP(), 'system', 1),
('Testing', UNIX_TIMESTAMP(), 'system', 1);