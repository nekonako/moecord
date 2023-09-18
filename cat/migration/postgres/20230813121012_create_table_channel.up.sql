CREATE TABLE IF NOT EXISTS "channel" (
    "id"            bytea           NOT NULL PRIMARY KEY,
    "server_id"     bytea           NOT NULL,
    "category_id"   bytea           NOT NULL,
    "name"          VARCHAR(255)    NOT NULL,
    "created_at"    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
