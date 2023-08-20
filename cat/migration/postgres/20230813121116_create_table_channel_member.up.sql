CREATE TABLE IF NOT EXISTS "channel_member" (
    "id"              bytea        NOT NULL PRIMARY KEY,
    "user_id"         bytea        NOT NULL,
    "channel_id"      bytea        NOT NULL,
    "created_at"      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);