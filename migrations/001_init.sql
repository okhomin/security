-- Write your migrate up statements here
BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE INDEX index_login ON users (login);

COMMIT;

---- create above / drop below ----

BEGIN;

DROP TABLE users;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
