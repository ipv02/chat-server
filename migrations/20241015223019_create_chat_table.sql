-- +goose Up
create table chat (
    id serial primary key,
    name text not null
);

create table chat_users (
    chat_id int not null,
    user_id int not null,
    primary key (chat_id, user_id)
);

create table messages (
    id serial primary key,
    user_id int not null,
    message text not null,
    created_at timestamp not null default now()
);

-- +goose Down
drop table chat;
drop table chat_users;
drop table messages;