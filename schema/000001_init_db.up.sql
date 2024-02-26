BEGIN;

CREATE TABLE users
(
    id serial primary key,
    email varchar(128) not null unique,
    password varchar(128) not null,
    refresh_token varchar(512) default '',
    verified bool not null default false,
    verif_code varchar(64) not null
);

CREATE TABLE files
(
    id serial primary key,
    name varchar(128) not null unique,
    location varchar(256) not null unique,
    created_at date not null default NOW(),
    updated_at date not null default NOW(),
    user_id int references users (id) on delete cascade
);

COMMIT;
