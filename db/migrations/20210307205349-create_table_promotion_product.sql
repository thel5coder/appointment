
-- +migrate Up
CREATE TABLE IF NOT EXISTS "promotion_products" (
  "id" char(36) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "product_id" char(36) NOT NULL,
  "promotion_id" char(36) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS "promotion_products";