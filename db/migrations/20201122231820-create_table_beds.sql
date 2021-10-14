-- +migrate Up
CREATE TABLE IF NOT EXISTS "beds"
(
    "id"         char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "bed_code"   varchar(5)           NOT NULL,
    "clinic_id"  char(36)             NOT NULL,
    "is_usable"  boolean,
    "created_at" timestamp            NOT NULL,
    "updated_at" timestamp            NOT NULL,
    "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "beds";