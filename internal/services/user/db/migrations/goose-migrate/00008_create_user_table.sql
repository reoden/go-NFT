-- +goose Up
-- +goose StatementBegin
ALTER TABLE "users" ADD COLUMN "certification" BOOLEAN DEFAULT FALSE;
ALTER TABLE "users" ADD COLUMN "real_name" VARCHAR(255) DEFAULT NULL;
ALTER TABLE "users" ADD COLUMN "id_card_no" VARCHAR(255) DEFAULT NULL;
ALTER TABLE "users" ADD COLUMN "user_role" VARCHAR(128) DEFAULT NULL;

COMMENT ON COLUMN users.certification IS '实名认证状态（TRUE或FALSE）';
COMMENT ON COLUMN users.real_name IS '真实姓名';
COMMENT ON COLUMN users.id_card_no IS '身份证号';
COMMENT ON COLUMN users.user_role IS '用户角色';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP COLUMN "certification";
ALTER TABLE "users" DROP COLUMN "real_name";
ALTER TABLE "users" DROP COLUMN "id_card_no";
ALTER TABLE "users" DROP COLUMN "user_role";
-- +goose StatementEnd
