DROP DATABASE IF EXISTS `message_service`;
CREATE DATABASE `message_service`;
USE `message_service`;

DROP TABLE IF EXISTS `email_smtp_configs`;
CREATE TABLE `email_smtp_configs`
(
    `id`         BIGINT(20)   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `host`       VARCHAR(255) NOT NULL COMMENT 'host',
    `port`       INT          NOT NULL COMMENT 'port',
    `username`   VARCHAR(255) NOT NULL COMMENT '用户名',
    `password`   VARCHAR(255) NOT NULL COMMENT '密码',
    `status`     TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '关闭还是开启',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT 'email smtp配置表';


DROP TABLE IF EXISTS `email_senders`;
CREATE TABLE `email_senders`
(
    `id`                   BIGINT(20)   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `email_smtp_config_id` BIGINT(20)   NOT NULL,
    `address`              VARCHAR(255) NOT NULL COMMENT 'email address',
    `created_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);
ALTER TABLE `email_senders`
    ADD CONSTRAINT address_uk UNIQUE (address);


DROP TABLE IF EXISTS `emails`;
CREATE TABLE `emails`
(
    `id`           BIGINT(20)    NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `uuid`         VARCHAR(255)  NOT NULL COMMENT '唯一id',
    `from`         VARCHAR(255)  NOT NULL COMMENT '发送者',
    `to`           VARCHAR(1024) NOT NULL COMMENT '接受者，逗号隔开多个接受者',
    `subject`      VARCHAR(1024) NOT NULL COMMENT '主题',
    `body`         TEXT          NOT NULL COMMENT '邮件内容',
    `status`       VARCHAR(255)  NOT NULL COMMENT '状态',
    `resend_count` INT           NOT NULL DEFAULT 0 COMMENT '重发次数',
    `resend_at`    DATETIME COMMENT '下次重发时间',
    `created_at`   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);
ALTER TABLE `emails`
    ADD CONSTRAINT uuid_uk UNIQUE (uuid);

DROP TABLE IF EXISTS `wechat_work_app_configs`;
CREATE TABLE `wechat_work_app_configs`
(
    `id`                          BIGINT(20)    NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `company_id`                  VARCHAR(255)  NOT NULL COMMENT '应用所属的企业id',
    `company_secret`              VARCHAR(1024) NOT NULL COMMENT '应用所属的企业id的secret',
    `agent_id`                    VARCHAR(255)  NOT NULL COMMENT '企微那边的app id',
    `msg_receiving_server_token`  VARCHAR(255)  NOT NULL COMMENT '接收消息服务器的token',
    `msg_receiving_server_aeskey` VARCHAR(255)  NOT NULL COMMENT '接收消息服务器的aeskey',
    `description`                 VARCHAR(1024) COMMENT '描述',
    `created_at`                  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`                  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);
ALTER TABLE `wechat_work_app_configs`
    ADD CONSTRAINT company_id_and_app_id_uk UNIQUE (company_id, agent_id);


DROP TABLE IF EXISTS `wechat_work_messages`;
CREATE TABLE `wechat_work_messages`
(
    `id`                 BIGINT(20)   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `uuid`               VARCHAR(255) NOT NULL COMMENT '唯一id',
    `from_app_config_id` BIGINT(20)   NOT NULL COMMENT '发送者, 使用哪个应用发送, 对应wechat_work_app_configs的id字段',
    `to`                 TEXT         NOT NULL COMMENT '接受者，逗号隔开多个接受者',
    `content`            TEXT         NOT NULL COMMENT '消息内容',
    `status`             VARCHAR(255) NOT NULL COMMENT '状态',
    `resend_count`       INT          NOT NULL DEFAULT 0 COMMENT '重发次数',
    `resend_at`          DATETIME COMMENT '下次重发时间',
    `created_at`         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);


DROP TABLE IF EXISTS `message_sending_logs`;
CREATE TABLE `message_sending_logs`
(
    `id`            BIGINT(20)   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `message_id`    BIGINT(20)   NOT NULL COMMENT '消息id',
    `is_success`    TINYINT(1)   NOT NULL COMMENT '是否成功',
    `failed_reason` text         NOT NULL COMMENT '失败原因',
    `type`          VARCHAR(255) NOT NULL COMMENT '消息类型，邮件还是微信消息',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);

INSERT INTO `email_smtp_configs`(host, port, username, password, status)
VALUES ('smtp.qq.com', 587, 'jeb.wang@foxmail.com', 'todo modify me', 1);


INSERT INTO `email_senders`(email_smtp_config_id, address)
VALUES (1, 'jeb.wang@foxmail.com');


INSERT INTO `wechat_work_app_configs`(company_id, company_secret, agent_id, msg_receiving_server_token,
                                      msg_receiving_server_aeskey,
                                      description)
VALUES ('todo modify me', 'todo modify me', '1000002', 'todo modify me',
        'todo modify me', '');

