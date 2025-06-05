-- Staging seed data for auth table - production-like but safe for testing
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES 
(1, 'admin', 'staging_admin_2024'),
(2, 'staging_user', 'staging_test_123');