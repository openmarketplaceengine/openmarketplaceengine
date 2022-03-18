create table if not exists vehicle
(
    id         text not null primary key,
    driver     text,
    model      text not null,
    make       text not null,
    color      text not null,
    plate      text not null,
    class      int4 not null,
    type       int4 not null,
    year       int4 not null,
    capacity   int4 not null,
    cargovol   int4,
    wheelchair int2 not null,
    childseats int2 not null,
    comment    text
);