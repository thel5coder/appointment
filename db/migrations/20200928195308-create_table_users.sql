-- +migrate Up
CREATE TABLE IF NOT EXISTS "users"
(
    "id"                 char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"          varchar(50)          NOT NULL,
    "email"              varchar(100)         NOT NULL,
    "mobile_phone"       varchar(20),
    "password"           varchar(128)         NOT NULL,
    "fcm_device_token"   varchar(255),
    "role_id"            char(36)             NOT NULL,
    "profile_picture_id" char(36)             NOT NULL,
    "is_active"          boolean,
    "activated_at"       timestamp,
    "created_at"         timestamp            NOT NULL,
    "updated_at"         timestamp            NOT NULL,
    "deleted_at"         timestamp
);
-- +migrate Down
DROP TABLE IF EXISTS "users";