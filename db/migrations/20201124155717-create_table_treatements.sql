-- +migrate Up
CREATE TABLE IF NOT EXISTS "master_treatments"
(
    "id"          char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"        varchar(255)         NOT NULL,
    "description" text,
    "duration"    int2,
    "photo_id"    char(36),
    "icon_id"     char(36),
    "created_at"  timestamp            NOT NULL,
    "updated_at"  timestamp            NOT NULL,
    "deleted_at"  timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "master_treatments";