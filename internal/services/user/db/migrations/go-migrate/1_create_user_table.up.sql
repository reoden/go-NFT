-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create "products" table
CREATE TABLE "public"."user" (
  "id" INT GENERATED ALWAYS AS IDENTITY,
  "user_id" uuid NOT NULL,
  "nickname" text NULL,
  "phone" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
