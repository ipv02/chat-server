-- +goose Up
create table chat (
    id serial primary key,
    name text not null
);

-- +goose Down
drop table  chat;
