-- Make password nullable for social logins
ALTER TABLE `blog_auth` MODIFY COLUMN `password` varchar(60) DEFAULT NULL COMMENT '密码';