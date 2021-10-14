-- +migrate Up
CREATE TABLE IF NOT EXISTS "doctor_week_date_schedules"
(
    "id"                   char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "week_day_schedule_id" char(36)             NOT NULL,
    "week_date"            timestamp            NOT NULL,
    "is_present"           bool,
    "created_at"           timestamp            NOT NULL,
    "updated_at"           timestamp            NOT NULL,
    "deleted_at"           timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "doctor_week_date_schedules";