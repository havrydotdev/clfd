BEGIN;

CREATE TABLE files
(
    id uuid primary key DEFAULT gen_random_uuid(),
    name varchar(128) not null unique,
    location varchar(256) not null unique
);

COMMIT;