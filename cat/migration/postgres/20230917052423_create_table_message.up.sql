CREATE TABLE IF NOT EXISTS "message"
    (
        "id"         bytea       NOT NULL PRIMARY KEY,
        "sender_id"  bytea       NOT NULL,
        "channel_id" bytea       NOT NULL,
        "content"    TEXT        NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ
    );
