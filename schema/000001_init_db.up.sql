BEGIN;

CREATE TABLE files
(
    id serial primary key,
    name varchar(128) not null unique,
    location varchar(256) not null unique,
    created_at date not null default NOW(),
    updated_at date not null default NOW()
);

COMMIT;