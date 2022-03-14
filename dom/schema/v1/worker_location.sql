create table if not exists worker_location
(
    recnum    bigserial   not null primary key,
    worker    text        not null,
    stamp     timestamptz not null,
    longitude float8      not null,
    latitude  float8      not null,
    speed     int4
);