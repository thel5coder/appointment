-- +migrate Up
CREATE TABLE IF NOT EXISTS "categories"
(
    "id"                 char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "slug"               varchar(255)         NOT NULL,
    "name"               varchar(255)         NOT NULL,
    "parent_id"          char(36),
    "is_active"          boolean,
    "file_icon_id"       char(36),
    "file_background_id" char(36),
    "created_at"         timestamp            NOT NULL,
    "updated_at"         timestamp            NOT NULL,
    "deleted_at"         timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "categories";