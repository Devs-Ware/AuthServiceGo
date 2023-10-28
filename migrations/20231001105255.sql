-- Create "users" table
CREATE TABLE "users" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "username" text NOT NULL,
  "policy_id" text NULL,
  "organization_id" text NULL,
  "totp_secret" text NOT NULL,
  "totp_url" text NOT NULL,
  "profile_id" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create "profiles" table
CREATE TABLE "profiles" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email_verified" boolean NULL DEFAULT false,
  "two_factor_enabled" boolean NULL DEFAULT false,
  "bio" text NULL,
  "user_id" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_profile" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "profiles" ("deleted_at");
