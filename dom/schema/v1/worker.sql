create table if not exists worker
(
    id         text not null primary key,
    status     int4 not null,
    rating     int4,
    jobs       int4,
    first_name text,
    last_name  text,
    vehicle    text,
    created    timestamptz,
    updated    timestamptz
);

create index if not exists idx_worker_status on worker(status);
