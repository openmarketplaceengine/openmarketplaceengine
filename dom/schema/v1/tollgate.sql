create table if not exists tollgate
(
    id         text not null primary key,
    name       text,
    b_boxes    jsonb,
    gate_line  jsonb,
    created_at timestamptz,
    updated_at timestamptz
);