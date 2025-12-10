-- +goose Up
CREATE TABLE "user_operate_stream" (
  "id" BIGSERIAL PRIMARY KEY,
  "gmt_create" timestamp with time zone DEFAULT NULL,
  "gmt_modified" timestamp with time zone DEFAULT NULL,
  "user_id" varchar(64) DEFAULT NULL,
  "type" varchar(64) DEFAULT NULL,
  "operate_time" timestamp with time zone DEFAULT NULL,
  "param" text,
  "extend_info" text,
  "deleted" integer DEFAULT NULL,
  "lock_version" integer DEFAULT NULL
);
COMMENT ON TABLE "user_operate_stream" IS '用户操作流水表';
COMMENT ON COLUMN "user_operate_stream"."id" IS '流水ID（自增主键）';
COMMENT ON COLUMN "user_operate_stream"."gmt_create" IS '创建时间';
COMMENT ON COLUMN "user_operate_stream"."gmt_modified" IS '最后更新时间';
COMMENT ON COLUMN "user_operate_stream"."user_id" IS '用户ID';
COMMENT ON COLUMN "user_operate_stream"."type" IS '操作类型';
COMMENT ON COLUMN "user_operate_stream"."operate_time" IS '操作时间';
COMMENT ON COLUMN "user_operate_stream"."param" IS '操作参数';
COMMENT ON COLUMN "user_operate_stream"."extend_info" IS '扩展字段';
COMMENT ON COLUMN "user_operate_stream"."deleted" IS '是否逻辑删除，0为未删除，非0为已删除';
COMMENT ON COLUMN "user_operate_stream"."lock_version" IS '乐观锁版本号';

-- +goose Down
DROP TABLE IF EXISTS "user_operate_stream";