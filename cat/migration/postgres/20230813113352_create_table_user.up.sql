CREATE TABLE IF NOT EXISTS "users" (
    "id"            bytea           NOT NULL PRIMARY KEY,
    "username"      VARCHAR(255)    NOT NULL,
    "email"         VARCHAR(255)    NOT NULL,
    "created_at"    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    CONSTRAINT user_email_unique    UNIQUE (email),
    CONSTRAINT user_username_unique UNIQUE (username)
);