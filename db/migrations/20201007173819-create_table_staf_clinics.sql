-- +migrate Up
CREATE TABLE IF NOT EXISTS "staff_clinics"
(
    "id"         char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "staff_id"   char(36)             NOT NULL,
    "clinic_id"  char(36)             NOT NULL,
    "created_at" timestamp            NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "staff_clinics";