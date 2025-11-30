-- Create "products" table
CREATE TABLE "public"."user" (
    "id" uuid NOT NULL,
  "user_id" text NOT NULL,
  "nickname" text NULL,
  "phone" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
