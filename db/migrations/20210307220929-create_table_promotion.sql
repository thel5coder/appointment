-- +migrate Up
CREATE TABLE IF NOT EXISTS "promotions"
(
    "id"                            char(36) PRIMARY KEY              NOT NULL DEFAULT (uuid_generate_v4()),
    "slug"                          varchar(255)                      NOT NULL,
    "name"                          varchar(100)                      NOT NULL,
    "customer_promotion_condition"  customer_promotion_condition_enum NOT NULL,
    "promotion_type"                promotion_type_enum               NOT NULL,
    "description"                   text,
    "start_date"                    timestamp,
    "end_date"                      timestamp,
    "foto_id"                       char(36)                          NOT NULL,
    "nominal_type"                  nominal_type_enum,
    "nominal_precentage"            int2,
    "nominal_amount"                float4,
    "birth_date_condition"          date,
    "sex_condition"                 sex_enum,
    "register_date_condition_start" date,
    "register_date_condition_end"   date,
    "created_at"                    timestamp                         NOT NULL,
    "updated_at"                    timestamp                         NOT NULL,
    "deleted_at"                    timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "promotions";