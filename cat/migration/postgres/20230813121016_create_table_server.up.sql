CREATE TABLE IF NOT EXISTS "server" (
    "id"              bytea           NOT NULL PRIMARY KEY,
    "owner_id"        bytea           NOT NULL,
    "name"            VARCHAR(255)    NOT NULL,
    "created_at"      TIMESTAMP       NOT NULL DEFAULT NOW(),
    "updated_at"      TIMESTAMP       NOT NULL DEFAULT NOW()
);