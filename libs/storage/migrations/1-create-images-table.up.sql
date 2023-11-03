create table images (
    name text primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);