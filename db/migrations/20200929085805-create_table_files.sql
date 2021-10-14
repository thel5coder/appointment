-- +migrate Up
CREATE TABLE IF NOT EXISTS "files"
(
    "id"         char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "path"       varchar(255),
    "created_at" timestamp            NOT NULL,
    "updated_at" timestamp            NOT NULL,
    "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "files";