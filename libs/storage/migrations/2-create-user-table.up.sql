create table users (
    id bigserial primary key,
    username  varchar unique,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

alter table images add column user_id bigint not null;

ALTER TABLE images ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);