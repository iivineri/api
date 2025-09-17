-- Migration up
-- [start][reset_passwords] ==================================================================================================
CREATE TABLE reset_passwords (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "email" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "expired_at" TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '1 day'
);
CREATE UNIQUE INDEX "reset_passwords_email_unique" ON "reset_passwords" ("email");
-- [end][reset_passwords] ====================================================================================================

-- [start][users] ============================================================================================================
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "nickname" character varying(32) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(128) NOT NULL,
  "enabled_2fa" boolean NOT NULL DEFAULT false,
  "secret_2fa" character varying(128),
  "date_of_birth" DATE NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ,
  "deleted_at" TIMESTAMPTZ,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "users_nickname_unique" ON "users" ("nickname") WHERE "deleted_at" IS NULL;
CREATE UNIQUE INDEX "users_email_unique" ON "users" ("email") WHERE "deleted_at" IS NULL;
CREATE INDEX "users_deleted_at_idx" ON "users" ("deleted_at");
-- [end][users] ============================================================================================================

-- [start][bans] ===========================================================================================================
CREATE TABLE bans (
    "id" serial PRIMARY KEY,
    "user_id" integer NOT NULL,
    "banned_by" integer NOT NULL,
    "reason" text NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT bans_user_id_fkey FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT,
    CONSTRAINT bans_banned_by_fkey FOREIGN KEY ("banned_by") REFERENCES users("id") ON DELETE RESTRICT
);
CREATE INDEX "bans_user_id_idx" ON bans("user_id");
CREATE INDEX "bans_banned_by_idx" ON bans("banned_by");
-- [end][bans] =============================================================================================================

-- [start][sessions] =======================================================================================================
CREATE TABLE "sessions" (
    "id" serial NOT NULL,
    "user_id" integer NOT NULL,
    "user_agent" character varying(255),
    "ip_address" character varying(45),
    -- "mime_type" character varying(100),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    PRIMARY KEY ("id"),
    CONSTRAINT "sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT
);
CREATE INDEX "sessions_user_id_idx" ON "sessions" ("user_id");
CREATE INDEX "sessions_deleted_at_idx" ON "sessions" ("deleted_at");
-- [end][sessions] =========================================================================================================

-- [start][user_2fa_secrets] ===============================================================================================
CREATE TABLE user_2fa_secrets (
    "id" serial NOT NULL,
    "user_id" integer NOT NULL,
    "hash" character varying(126) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    PRIMARY KEY ("id"),
    CONSTRAINT "user_2fa_secrets_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE INDEX "user_2fa_secrets_user_id_idx" ON user_2fa_secrets ("user_id");
CREATE INDEX "user_2fa_secrets_deleted_at_idx" ON user_2fa_secrets ("deleted_at");
-- [end][user_2fa_secrets] =================================================================================================

-- [start][images] =========================================================================================================
CREATE TABLE "images" (
    "id" serial NOT NULL,
    "user_id" integer NOT NULL,
    "thumber_id" character varying(64) NOT NULL,
    "original_name" character varying(255),
    "mime_type" character varying(100),
    "size" bigint NOT NULL DEFAULT 0,
    "width" integer NOT NULL DEFAULT 0,
    "height" integer NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    PRIMARY KEY ("id"),
    CONSTRAINT "images_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE INDEX "images_user_id_idx" ON "images" ("user_id");
CREATE INDEX "images_deleted_at_idx" ON "images" ("deleted_at");
-- [end][images] ===========================================================================================================

-- [start][accepted_images] ================================================================================================
CREATE TABLE "accepted_images" (
  "id" serial NOT NULL,
  "image_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "approved_by" integer NOT NULL,
  "nsfw" boolean NOT NULL DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("id"),
  CONSTRAINT "accepted_images_image_id_fkey" FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE RESTRICT,
  CONSTRAINT "accepted_images_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT,
  CONSTRAINT "accepted_images_approved_by_fkey" FOREIGN KEY ("approved_by") REFERENCES "users" ("id") ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "accepted_images_image_id_unique" ON "accepted_images" ("image_id");
-- [end][accepted_images] ==================================================================================================

-- [start][rejected_images] ================================================================================================
CREATE TABLE "rejected_images" (
  "id" serial NOT NULL,
  "image_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "rejected_by" integer NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("id"),
  CONSTRAINT "rejected_images_image_id_fkey" FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE,
  CONSTRAINT "rejected_images_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  CONSTRAINT "rejected_images_rejected_by_fkey" FOREIGN KEY ("rejected_by") REFERENCES "users" ("id") ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "rejected_images_image_id_unique" ON "rejected_images" ("image_id");
-- [end][rejected_images] ==================================================================================================

-- [start][activities] =====================================================================================================
CREATE TABLE activities (
    "id" BIGSERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "table_name" VARCHAR(255) NOT NULL,
    "record_id" VARCHAR(255) NOT NULL,
    "action" VARCHAR(50) NOT NULL,
    "changes" JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id"),
    CONSTRAINT "activities_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT
);
CREATE INDEX "activities_table_name_record_id_idx" ON activities ("table_name", "record_id");
CREATE INDEX "activities_user_id_idx" ON activities ("user_id");
CREATE INDEX "activities_created_at_idx" ON activities ("created_at");
-- [end][activities] =======================================================================================================
