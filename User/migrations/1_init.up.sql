CREATE TABLE IF NOT EXISTS users
(
    name text primary key,
    messages_count int check (messages_count > 0),
    created timestamp with time zone
);
