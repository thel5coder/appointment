
-- +migrate Up
CREATE TYPE "sex_enum" AS ENUM (
  'male',
  'female'
  'other'
);

CREATE TYPE "religion_enum" AS ENUM (
  'islam',
  'katholik',
  'protestan',
  'hindhu',
  'budha',
  'other'
);

CREATE TABLE "customers" (
  "id" char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "name" varchar(50) NOT NULL,
  "sex" sex_enum NOT NULL,
  "address" text,
  "birth_date" date,
  "marital_status" varchar(50),
  "phone_number" varchar(20),
  "mobile_phone_1" varchar(20),
  "mobile_phone_2" varchar(20),
  "religion" religion_enum,
  "education" varchar(50),
  "hobby" varchar(50),
  "profession" varchar(30),
  "reference" varchar(50),
  "notes" text,
  "city_id" char(36) NOT NULL,
  "user_id" char(36) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE "customers";
