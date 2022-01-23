-- auto-generated definition
create table users
(
    id         uuid not null
        constraint users_pk
            primary key,
    name       text not null,
    occupation text not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

alter table users
    owner to customuser;

