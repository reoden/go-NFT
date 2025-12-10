-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user_operate_stream" (
  "id" BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID（自增主键）',
  "gmt_create" datetime DEFAULT NULL COMMENT '创建时间',
  "gmt_modified" datetime DEFAULT NULL COMMENT '最后更新时间',
  "user_id" varchar(64) DEFAULT NULL COMMENT '用户ID',
  "type" varchar(64) DEFAULT NULL COMMENT '操作类型',
  "operate_time" datetime DEFAULT NULL COMMENT '操作时间',
  "param" text COMMENT '操作参数',
  "extend_info" text COMMENT '扩展字段',
  "deleted" int DEFAULT NULL COMMENT '是否逻辑删除，0为未删除，非0为已删除',
  "lock_version" int DEFAULT NULL COMMENT '乐观锁版本号',
  PRIMARY KEY ("id")
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户操作流水表';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_operate_stream";
-- +goose StatementEnd