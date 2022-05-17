create table if not exists jobimp
(
    id           text not null primary key,
    worker_id    text not null,
    created      timestamptz,
    updated      timestamptz,
    state        text,
    pickup_date  timestamptz,
    pickup_addr  text,
    pickup_lat   float8,
    pickup_lon   float8,
    dropoff_addr text,
    dropoff_lat  float8,
    dropoff_lon  float8,
    trip_type    text,
    category     text
);
