
-- +migrate Up
CREATE TABLE IF NOT EXISTS "doctor_week_day_schedules" (
                                             "id" char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
                                             "doctor_id" char(36) NOT NULL,
                                             "clinic_id" char(36) NOT NULL,
                                             "day" varchar(50) NOT NULL,
                                             "schedule_date" date,
                                             "is_present" boolean,
                                             "created_at" timestamp NOT NULL,
                                             "updated_at" timestamp NOT NULL,
                                             "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "doctor_week_day_schedules";