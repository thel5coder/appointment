
-- +migrate Up
create extension if not exists "uuid-ossp";
-- +migrate Down
drop extension if exists "uuid-ossp";