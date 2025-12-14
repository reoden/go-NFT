-- +goose Up
-- +goose StatementBegin
ALTER TABLE "users" ADD COLUMN "state" VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP COLUMN "state";
-- +goose StatementEnd
