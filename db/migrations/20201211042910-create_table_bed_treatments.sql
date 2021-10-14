-- +migrate Up
CREATE TABLE IF NOT EXISTS "bed_treatments"
(
    "id"           char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "bed_id"       char(36)             NOT NULL,
    "treatment_id" char(36)             NOT NULL,
    "created_at"   timestamp            NOT NULL,
    "updated_at"   timestamp            NOT NULL,
    "deleted_at"   timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "bed_treatments";