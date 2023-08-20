CREATE TABLE IF NOT EXISTS "server_member" (
    "id"              bytea        NOT NULL PRIMARY KEY,
    "user_id"         bytea        NOT NULL,
    "server_id"       bytea        NOT NULL,
    "created_at"      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);