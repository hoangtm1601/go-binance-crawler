drop index public.idx_symbol_interval_start;

create unique index idx_symbol_interval_start
    on public.candles (symbol, interval, start);