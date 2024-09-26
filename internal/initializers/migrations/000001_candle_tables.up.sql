CREATE TABLE candles (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    symbol VARCHAR(255) NOT NULL,
    interval VARCHAR(255) NOT NULL,
    start BIGINT NOT NULL,
    "end" BIGINT NOT NULL,
    last_end BIGINT,
    op FLOAT8,
    hi FLOAT8,
    lo FLOAT8,
    cl FLOAT8,
    bv FLOAT8,
    qv FLOAT8,
    tbv FLOAT8,
    tqv FLOAT8,
    cnt BIGINT
);

CREATE INDEX idx_symbol_interval_start ON candles (symbol, interval, start);
