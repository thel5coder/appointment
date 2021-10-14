-- +migrate Up

CREATE TABLE IF NOT EXISTS "doctor_weekday_schedule_work_times"
(
    "id"                   char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "week_day_schedule_id" char(36)             NOT NULL,
    "work_time"            timestamp            NOT NULL,
    "created_at"           timestamp            NOT NULL,
    "updated_at"           timestamp            NOT NULL,
    "deleted_at"           timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "doctor_weekday_schedule_work_times";