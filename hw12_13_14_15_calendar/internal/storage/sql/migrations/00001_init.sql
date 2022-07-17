-- +goose Up
create table events
(
    id            varchar(36) primary key,
    title         varchar     not null,
    start_at      timestamp   not null,
    end_at        timestamp   not null,
    description   text        null,
    owner_id      varchar(36) not null,
    notify_before interval    null
);

create index if not exists events_owner_date_range_idx on events(owner_id, start_at, end_at);

-- +goose Down
drop table events;