CREATE TABLE IF NOT EXISTS ping_results (
    ip VARCHAR(15) NOT NULL,
    ping_time TIMESTAMP NOT NULL,
    last_success TIMESTAMP NULL
);