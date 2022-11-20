CREATE TABLE IF NOT EXISTS tg_query_executor
(
    id                SERIAL PRIMARY KEY,
    bot_token         TEXT        NOT NULL,
    bot_chat_id       BIGINT      NOT NULL,
    chat_describe     TEXT        NOT NULL,
    sql_query         TEXT        NOT NULL,
    schedule_cron     VARCHAR(20) NOT NULL,
    last_execution_ts TIMESTAMP
);