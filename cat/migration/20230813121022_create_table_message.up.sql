CREATE TABLE IF NOT EXISTS "message" (
    "id"          bytea        NOT NULL PRIMARY KEY,
    "sender_id"   bytea         NOT NULL,
    "data"        TEXT         NOT NULL,
    "type"        VARCHAR(255) NOT NULL,
    "created_at"  TIMESTAMP    NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMP    NOT NULL DEFAULT NOW()
)