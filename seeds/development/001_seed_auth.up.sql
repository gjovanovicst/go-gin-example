-- Development seed data for auth table with bcrypt hashed passwords
-- admin123 -> $2a$10$Pca7n7ahfOgCXp8MAm51xeVl8Hv7uuopr9w1elAd7Y6qEODer8CFS
-- test123 -> $2a$10$YVIzGkB9gI2M5s/2AqAkBO/Z8CdWjnnB8CLTrThvCIo7fjEvXEina
-- dev123 -> $2a$10$buNpVLQB8f.3ktB9OLtMjeKxv9zOlbDKB3E5N8t91FvdBgL4LGope
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES
(1, 'admin', '$2a$10$Pca7n7ahfOgCXp8MAm51xeVl8Hv7uuopr9w1elAd7Y6qEODer8CFS'),
(2, 'testuser', '$2a$10$YVIzGkB9gI2M5s/2AqAkBO/Z8CdWjnnB8CLTrThvCIo7fjEvXEina'),
(3, 'developer', '$2a$10$buNpVLQB8f.3ktB9OLtMjeKxv9zOlbDKB3E5N8t91FvdBgL4LGope');