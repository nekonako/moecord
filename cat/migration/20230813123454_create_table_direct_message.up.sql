CREATE TABLE IF NOT EXISTS "direct_message" (
    "id"              bytea       NOT NULL PRIMARY KEY,
    "message_id"      bytea       NOT NULL,
    "receiver_id"     bytea       NOT NULL,
    "created_at"      TIMESTAMP   NOT NULL DEFAULT NOW()
)