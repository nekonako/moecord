CREATE TABLE IF NOT EXISTS "servers" (
    "id"              bytea           NOT NULL PRIMARY KEY,
    "owner_id"        bytea           NOT NULL,
    "name"            VARCHAR(255)    NOT NULL,
    "direct_message"  BOOLEAN         NOT NULL DEFAULT false,
    "created_at"      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    "updated_at"      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);