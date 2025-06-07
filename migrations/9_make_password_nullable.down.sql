-- Restore password as NOT NULL
ALTER TABLE `blog_auth` MODIFY COLUMN `password` varchar(50) DEFAULT '' COMMENT '密码';