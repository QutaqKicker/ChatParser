CREATE TABLE IF NOT EXISTS users
(
    id text primary key,
    name text not null unique,
    messages_count int check (messages_count > 0),
    created timestamp with time zone
);
