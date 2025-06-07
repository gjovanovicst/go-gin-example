-- Production seed data for auth table - minimal essential data with bcrypt hashed passwords
-- admin123 -> $2a$10$Pca7n7ahfOgCXp8MAm51xeVl8Hv7uuopr9w1elAd7Y6qEODer8CFS
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES
(1, 'admin', '$2a$10$Pca7n7ahfOgCXp8MAm51xeVl8Hv7uuopr9w1elAd7Y6qEODer8CFS');