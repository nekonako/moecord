CREATE TABLE IF NOT EXISTS channel_category (
    "id"            bytea        NOT NULL PRIMARY KEY,
    "server_id"     bytea        NOT NULL,
    "name"          VARCHAR(255) NOT NULL,
    "created_at"    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
