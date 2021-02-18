-- 用户表
CREATE TABLE `gin_user`
(
    `user_id`     int(20) NOT NULL AUTO_INCREMENT,
    `username`    varchar(50) NOT NULL COMMENT '用户名',
    `password`    varchar(50) NOT NULL COMMENT '密码',
    `salt`        varchar(16) NOT NULL COMMENT '盐',
    `email`       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `create_time` varchar(50) NOT NULL COMMENT '创建时间',
    `update_time` varchar(50)          DEFAULT NULL COMMENT '更新时间',
    `is_admin`    int(11) NOT NULL DEFAULT '0' COMMENT '是否是 admin 用户',
    `status`      int(11) NOT NULL DEFAULT '1' COMMENT '启用状态',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `UQE_user_name` (`username`),
    UNIQUE KEY `UQE_user_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户表';

INSERT INTO test.vcd_user
(username, password, salt, email, create_time, update_time, is_admin, status)
VALUES ('admin', 'f3e251e3242469a361cdf2653a75f70e', 'JWK1tU', '123@123.com', '2020-06-13 12:22:18',
        '2020-06-20 14:49:36', 1, 1);

CREATE TABLE `gin_user_token`
(
    `user_id`     int(20) NOT NULL COMMENT '用户 id',
    `token`       varchar(200) NOT NULL COMMENT 'token',
    `expire_time` varchar(50) DEFAULT '' COMMENT '过期时间',
    `update_time` varchar(50) DEFAULT '' COMMENT '更新时间',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `UQE_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户Token';

CREATE TABLE `gin_host`
(
    `host_id`     int(20) NOT NULL AUTO_INCREMENT COMMENT '主机 id',
    `host_alias`  varchar(100) NOT NULL DEFAULT '' COMMENT '主机别名',
    `host_name`   varchar(100) NOT NULL COMMENT '主机名',
    `host_port`   int(11) NOT NULL DEFAULT '5921' COMMENT '端口号',
    `remark`      varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
    `create_time` varchar(50)  NOT NULL COMMENT '创建时间',
    `update_time` varchar(50)           DEFAULT '' COMMENT '更新时间',
    PRIMARY KEY (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机管理';

CREATE TABLE `gin_task`
(
    `task_id`            int(20) NOT NULL AUTO_INCREMENT COMMENT '任务 id',
    `task_name`          varchar(100) NOT NULL COMMENT '任务名称',
    `task_tag`           varchar(100) NOT NULL DEFAULT '' COMMENT '任务标签',
    `task_level`         int(11) NOT NULL DEFAULT '1' COMMENT '任务等级 1: 主任务 2: 依赖任务',
    `dependency_task_id` varchar(100) NOT NULL DEFAULT '' COMMENT '依赖任务ID,多个ID逗号分隔',
    `dependency_status`  int(11) NOT NULL DEFAULT '1' COMMENT '依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖',
    `cron_expression`    varchar(100) NOT NULL COMMENT 'crontab 表达式',
    `protocol`           int(11) NOT NULL COMMENT '协议 1:http 2:系统命令',
    `http_method`        int(11) NOT NULL DEFAULT '1' COMMENT 'http请求方法',
    `command`            varchar(300) NOT NULL COMMENT 'URL地址或shell命令',
    `timeout`            int(11) NOT NULL DEFAULT '0' COMMENT '任务执行超时时间(单位秒),0不限制',
    `is_multi_instance`  int(11) NOT NULL DEFAULT '1' COMMENT '是否允许多实例运行',
    `retry_times`        int(11) NOT NULL DEFAULT '0' COMMENT '重试次数',
    `retry_interval`     int(11) NOT NULL DEFAULT '0' COMMENT '重试间隔时间',
    `notify_status`      int(11) NOT NULL DEFAULT '1' COMMENT '任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知',
    `notify_type`        int(11) NOT NULL DEFAULT '0' COMMENT '通知类型 1: 邮件 2: slack 3: webhook',
    `notify_receiver_id` varchar(300) NOT NULL DEFAULT '' COMMENT '通知接受者ID, setting表主键ID，多个ID逗号分隔',
    `notify_keyword`     varchar(200) NOT NULL DEFAULT '' COMMENT '通知关键词',
    `task_remark`        varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
    `task_status`        int(11) NOT NULL DEFAULT '0' COMMENT '状态 1:正常 0:停止',
    `create_time`        varchar(50)  NOT NULL COMMENT '创建时间',
    `update_time`        varchar(50)           DEFAULT '' COMMENT '更新时间',
    `delete_time`        varchar(50)           DEFAULT '' COMMENT '删除时间',
    PRIMARY KEY (`task_id`),
    KEY                  `IDX_task_level` (`task_level`),
    KEY                  `IDX_task_protocol` (`protocol`),
    KEY                  `IDX_task_status` (`task_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务管理';

CREATE TABLE `gin_task_host`
(
    `id`      int(20) NOT NULL AUTO_INCREMENT,
    `task_id` int(20) NOT NULL COMMENT '任务 id',
    `host_id` int(20) NOT NULL COMMENT '主机 id',
    PRIMARY KEY (`id`),
    KEY       `IDX_task_host_task_id` (`task_id`),
    KEY       `IDX_task_host_host_id` (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务-主机管理';

CREATE TABLE `gin_task_log`
(
    `id`              bigint(20) NOT NULL AUTO_INCREMENT,
    `task_id`         int(20) NOT NULL DEFAULT '0',
    `task_name`       varchar(100) NOT NULL,
    `cron_expression` varchar(100) NOT NULL,
    `protocol`        int(11) NOT NULL,
    `command`         varchar(300) NOT NULL,
    `timeout`         int(11) NOT NULL DEFAULT '0',
    `retry_times`     int(11) NOT NULL DEFAULT '0',
    `host_name`       varchar(200) NOT NULL DEFAULT '',
    `start_time`      datetime              DEFAULT NULL,
    `end_time`        datetime              DEFAULT NULL,
    `status`          int(11) NOT NULL DEFAULT '1',
    `result`          text         NOT NULL,
    PRIMARY KEY (`id`),
    KEY               `IDX_task_log_task_id` (`task_id`),
    KEY               `IDX_task_log_protocol` (`protocol`),
    KEY               `IDX_task_log_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务-日志管理';

CREATE TABLE `gin_login_log`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `username`    varchar(100) NOT NULL COMMENT '用户名',
    `ip`          varchar(64)  NOT NULL COMMENT 'IP地址',
    `create_time` varchar(50)  NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登陆日志管理';

