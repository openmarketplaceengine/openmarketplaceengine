create table if not exists tollgate_crossing
(
    id          text        not null primary key,
    tollgate_id text        not null,
    driver_id   text        not null,
    crossing    jsonb       not null,
    created     timestamptz not null
);
