CREATE TABLE IF NOT EXISTS chats
(
    id integer primary key generated always as identity,
    name text not null unique,
    created timestamp with time zone
);

CREATE TABLE IF NOT EXISTS users
(
    id integer primary key generated always as identity,
    name text not null unique,
    created timestamp with time zone
);

CREATE TABLE IF NOT EXISTS messages
(
    id uuid primary key,
    chat_id integer,
    user_id integer,
    text text,
    created timestamp with time zone,
    constraint fk_chat
        foreign key(chat_id)
            references chats(id),
    constraint fk_user
        foreign key(user_id)
            references users(id)
);

CREATE INDEX IF NOT EXISTS idx_created on messages(created);
