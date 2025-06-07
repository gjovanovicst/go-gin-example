-- Make username nullable for social logins
ALTER TABLE `blog_auth` MODIFY COLUMN `username` varchar(50) DEFAULT NULL COMMENT '账号';