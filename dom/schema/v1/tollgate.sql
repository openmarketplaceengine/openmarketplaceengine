create table if not exists tollgate
(
    id        text        not null primary key,
    name      text        not null,
    b_boxes   jsonb,
    gate_line jsonb,
    created   timestamptz not null,
    updated   timestamptz
);