-- Staging seed data for auth table - production-like but safe for testing with bcrypt hashed passwords
-- staging_admin_2024 -> $2a$10$7xa5tDm7j6kueF6XRJdxSO.milMreBsEA6Xs1Iva0GFPuNzGggvzO
-- staging_test_123 -> $2a$10$0AdUqzbr8kvNH5O/LqIfOuafPWkVnrQviMj3KAcglBykjhejjJHqa
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES
(1, 'admin', '$2a$10$7xa5tDm7j6kueF6XRJdxSO.milMreBsEA6Xs1Iva0GFPuNzGggvzO'),
(2, 'staging_user', '$2a$10$0AdUqzbr8kvNH5O/LqIfOuafPWkVnrQviMj3KAcglBykjhejjJHqa');