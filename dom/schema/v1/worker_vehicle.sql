create table if not exists worker_vehicle
(
    worker  text not null,
    vehicle text not null,
    PRIMARY KEY (worker, vehicle)
);