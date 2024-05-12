-- +goose Up
create table if not exists public.sessions
(
    id          uuid         not null
        primary key,
    user_id     uuid         not null,
    fingerprint varchar(255) not null,
    expires_at  timestamp(0),
    created_at  timestamp(0),
    updated_at  timestamp(0)
);

-- +goose Down
drop table if exists public.sessions;

