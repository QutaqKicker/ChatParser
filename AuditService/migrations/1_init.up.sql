CREATE TABLE IF NOT EXISTS logs
(
    id uuid primary key DEFAULT gen_random_uuid(),
    service_name text not null,
    "type" int not null,
    message text not null,
    created timestamp with time zone default now()
);
