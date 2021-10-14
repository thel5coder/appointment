-- +migrate Up
CREATE TABLE IF NOT EXISTS "doctor_treatments"
(
    "id"           char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "doctor_id"    char(36)             NOT NULL,
    "treatment_id" char(36)             NOT NULL,
    "created_at"   timestamp            NOT NULL,
    "updated_at"   timestamp            NOT NULL,
    "deleted_at"   timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "doctor_treatments";