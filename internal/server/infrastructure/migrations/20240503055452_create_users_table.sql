-- +goose Up
create table if not exists public.users
(
    id                uuid         not null
        primary key,
    email             varchar(255) not null
        constraint users_email_unique
            unique,
    password          varchar(255) not null,
    verification_code varchar(255) not null,
    version           bigint       not null,
    verified_at       timestamp(0),
    created_at        timestamp(0),
    updated_at        timestamp(0)
);

create index if not exists users_email_index
    on public.users (email);

-- +goose Down
drop index if exists users_email_index;
drop table if exists public.users;
