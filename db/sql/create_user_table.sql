CREATE TABLE table_user (
    id INT(11) NOT NULL AUTO_INCREMENT,
    user_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
    user_pwd VARCHAR(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
    email VARCHAR(64) DEFAULT '' COMMENT '邮箱',
    phone VARCHAR(128) DEFAULT '' COMMENT '手机号',
    email_validated TINYINT(1) DEFAULT 0 COMMENT '邮箱是否已验证',
    phone_validated TINYINT(1) DEFAULT 0 COMMENT '手机号是否已验证',
    signup_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
    last_active DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
    profile TEXT COMMENT '用户属性',
    status INT(11) NOT NULL DEFAULT 0 COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
    PRIMARY KEY (id),
    UNIQUE KEY idx_phone (phone),
    KEY idx_status (status)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;
