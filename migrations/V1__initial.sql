CREATE SCHEMA trades;

CREATE TABLE trades.klines (
    pair varchar NOT NULL,
    time_frame VARCHAR NOT NULL,
    open_price DECIMAL NOT NULL,
    high_price DECIMAL NOT NULL,
    low_price DECIMAL NOT NULL,
    close_price DECIMAL NOT NULL,
    utc_begin BIGINT NOT NULL,
    utc_end BIGINT NOT NULL,
    buy_base DECIMAL NOT NULL DEFAULT 0.0,
    sell_base DECIMAL NOT NULL DEFAULT 0.0,
    buy_quote DECIMAL NOT NULL DEFAULT 0.0,
    sell_quote DECIMAL NOT NULL DEFAULT 0.0
);

CREATE TABLE trades.recent_trades (
    tid varchar NOT NULL,
    pair varchar NOT NULL,
    price varchar NOT NULL,
    amount varchar NOT NULL,
    side varchar NOT NULL,
    ts BIGINT NOT NULL
);
