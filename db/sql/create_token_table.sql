CREATE TABLE `table_token` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_name` VARCHAR(40) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` CHAR(40) NOT NULL DEFAULT '' COMMENT '用户登录 token',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
