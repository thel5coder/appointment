-- +migrate Up
CREATE TABLE IF NOT EXISTS "treatment_procedures"
(
    "id"           char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "procedure_id" char(36)             NOT NULL,
    "treatment_id" char(36)             NOT NULL,
    "created_at"   timestamp            NOT NULL,
    "updated_at"   timestamp            NOT NULL,
    "deleted_at"   timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "treatment_procedures";