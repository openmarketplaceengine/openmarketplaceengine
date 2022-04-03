create table if not exists tollgate_crossing
(
    id          text not null primary key,
    tollgate_id text,
    driver_id   text,
    crossing    jsonb,
    created_at  timestamptz,
    constraint fk_tollgate foreign key (tollgate_id) references tollgate (id)
);
