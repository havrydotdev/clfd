BEGIN;

CREATE TABLE users
(
    id serial primary key,
    email varchar(128) not null unique,
    refresh_token varchar(64) not null unique,
    verified bool not null default false,
    verif_code varchar(64)
);

CREATE TABLE files
(
    id serial primary key,
    name varchar(128) not null unique,
    location varchar(256) not null unique,
    created_at date not null default NOW(),
    updated_at date not null default NOW()
);

COMMIT;