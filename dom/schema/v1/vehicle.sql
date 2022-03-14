create table if not exists vehicle
(
    id         text not null primary key,
    driver     text,
    model      text,
    make       text,
    color      text,
    plate      text,
    class      int4,
    type       int4,
    year       int4,
    capacity   int4,
    cargovol   int4,
    wheelchair int2,
    childseats int2,
    comment    text
);