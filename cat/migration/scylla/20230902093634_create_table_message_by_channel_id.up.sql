CREATE MATERIALIZED VIEW IF NOT EXISTS messages_by_channel_id as 
    SELECT * FROM  messages WHERE channel_id IS NOT NULL 
    PRIMARY KEY(channel_id, id) WITH CLUSTERING ORDER BY (id ASC);