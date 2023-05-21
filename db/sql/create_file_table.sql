CREATE TABLE `table_file` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` CHAR(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` BIGINT(20) DEFAULT 0 COMMENT '文件大小',
  `file_addr` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '文件存储路径',
  `create_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` INT(11) DEFAULT 0 COMMENT '状态(可用/禁用/已删除/其它)',
  `ext1` INT(11) DEFAULT 0 COMMENT '备用字段',
  `ext2` TEXT COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
