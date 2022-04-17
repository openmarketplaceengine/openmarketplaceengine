create table if not exists worker_status_enum
(
    id   int4 not null primary key,
    name text not null
);

BEGIN;
delete
from worker_status_enum;
insert into worker_status_enum
values (0, 'offline'),
       (1, 'ready'),
       (2, 'onjob'),
       (3, 'paused'),
       (4, 'disabled');
COMMIT;
