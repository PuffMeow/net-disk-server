CREATE TABLE table_user_file (
  id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  user_name varchar(64) NOT NULL COMMENT '用户名',
  file_sha1 varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  file_size bigint(20) DEFAULT '0' COMMENT '文件大小',
  file_name varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  upload_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  last_update datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  status int(11) NOT NULL DEFAULT '0' COMMENT '文件状态 (0 正常, 1 已删除, 2 禁用)',
  KEY idx_status (status),
  KEY idx_user_name (user_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;