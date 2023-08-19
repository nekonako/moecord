CREATE TABLE IF NOT EXISTS "channel" (
    "id"            bytea           NOT NULL PRIMARY KEY,
    "name"          VARCHAR(255)    NOT NULL,
    "avatar"        VARCHAR(255)    NOT NULL DEFAULT '',
    "created_at"    TIMESTAMP       NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP       NOT NULL DEFAULT NOW()
);