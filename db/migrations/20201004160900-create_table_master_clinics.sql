-- +migrate Up
CREATE TABLE IF NOT EXISTS "master_clinics"
(
    "id"           char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"         varchar(100)         NOT NULL,
    "address"      text,
    "pic_name"     varchar(50),
    "phone_number" varchar(20),
    "email"        varchar(50),
    "created_at"   timestamp            NOT NULL,
    "updated_at"   timestamp            NOT NULL,
    "deleted_at"   timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "master_clinics";