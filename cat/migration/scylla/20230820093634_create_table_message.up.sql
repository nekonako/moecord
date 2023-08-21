CREATE TABLE IF NOT EXISTS messages (
    id BLOB,
    channel_id BLOB,
    sender_id BLOB,
    content TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (channel_id, id)
) WITH CLUSTERING ORDER BY (id DESC);