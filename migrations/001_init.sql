-- Write your migrate up statements here
BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE files
(
    id      uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name    TEXT   NOT NULL UNIQUE,
    content TEXT   NOT NULL,
    groups  uuid[] NOT NULL
);

CREATE TABLE groups
(
    id    uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  TEXT    NOT NULL UNIQUE,
    read  BOOLEAN NOT NULL,
    write BOOLEAN NOT NULL,
    users uuid[]  NOT NULL
);

CREATE INDEX index_groups_name ON groups (name);
CREATE INDEX index_files_name ON files (name);
CREATE INDEX index_users_login ON users (login);

COMMIT;

---- create above / drop below ----

BEGIN;

DROP TABLE users;
DROP TABLE files;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
