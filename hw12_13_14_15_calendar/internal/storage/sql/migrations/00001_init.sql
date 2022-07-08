-- +goose Up
create table events
(
    id          varchar(16) primary key,
    title       varchar     not null,
    start_at    timestamp   not null,
    end_at      timestamp   not null,
    description text        null,
    owner_id    varchar(16) not null,
    notify_at   timestamp   null
);

create index if not exists events_owner_date_range_idx on events using btree (owner_id, start_at, end_at);
create index if not exists events_notify_idx on events using btree (notify_at) where notify_at is not null;

-- +goose Down
drop table events;