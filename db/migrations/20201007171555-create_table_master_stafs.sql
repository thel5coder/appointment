-- +migrate Up
CREATE TABLE iF NOT EXISTS "master_staffs"
(
    "id"          char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"        varchar(50)          NOT NULL,
    "description" text,
    "user_id"     char(36)             NOT NULL,
    "created_at"  timestamp            NOT NULL,
    "updated_at"  timestamp            NOT NULL,
    "deleted_at"  timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "master_staffs";
