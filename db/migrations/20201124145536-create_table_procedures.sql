
-- +migrate Up
CREATE TABLE IF NOT EXISTS "procedures" (
                              "id" char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
                              "name" varchar(255) NOT NULL,
                              "duration" int2,
                              "created_at" timestamp NOT NULL,
                              "updated_at" timestamp NOT NULL,
                              "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "procedures";