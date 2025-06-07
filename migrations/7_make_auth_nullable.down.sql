-- Restore username and password as NOT NULL
ALTER TABLE `blog_auth` MODIFY COLUMN `username` varchar(50) DEFAULT '' COMMENT '账号';
ALTER TABLE `blog_auth` MODIFY COLUMN `password` varchar(50) DEFAULT '' COMMENT '密码';