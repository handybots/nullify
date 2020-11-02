CREATE DATABASE IF NOT EXISTS nullifybot;

CREATE TABLE IF NOT EXISTS nullifybot.logs (
    date Date,
    time DateTime,
    level String,
    message String,
    event String,
    user_id UInt32
) ENGINE = MergeTree(date, (level, event, user_id), 8192);
