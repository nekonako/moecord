CREATE TABLE IF NOT EXISTS "user" (
    "id"            bytea           NOT NULL PRIMARY KEY,
    "username"      VARCHAR(255)    NOT NULL,
    "email"         VARCHAR(255)    NOT NULL,
    "password"      VARCHAR(255)    NOT NULL,
    "created_at"    TIMESTAMP       NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP       NOT NULL DEFAULT NOW()
);