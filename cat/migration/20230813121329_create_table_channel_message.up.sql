CREATE TABLE IF NOT EXISTS "channel_message" (
    "id"              bytea       NOT NULL PRIMARY KEY,
    "channel_id"      bytea       NOT NULL,
    "message_id"      bytea       NOT NULL,
    "created_at"      TIMESTAMP   NOT NULL DEFAULT NOW()
)