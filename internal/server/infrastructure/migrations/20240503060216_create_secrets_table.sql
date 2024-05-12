-- +goose Up
create table if not exists public.secrets
(
    id         uuid         not null
        primary key,
    user_id    uuid         not null,
    type       smallint     not null,
    name       varchar(255) not null,
    data       jsonb,
    notes      varchar(255),
    created_at timestamp(0),
    updated_at timestamp(0),
    deleted_at timestamp(0)
);
create unique index if not exists idx_id_user_id
    on public.secrets (id, user_id);

-- +goose Down
drop index if exists idx_id_user_id;
drop table if exists public.secrets;
