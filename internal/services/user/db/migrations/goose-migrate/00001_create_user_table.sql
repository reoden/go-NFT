-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
                        "id" SERIAL PRIMARY KEY,
                        "user_id" uuid NOT NULL,
                        "nickname" text NULL,
                        "phone" text NULL,
                        "created_at" timestamptz NULL,
                        "updated_at" timestamptz NULL,
                        "deleted_at" timestamptz NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
