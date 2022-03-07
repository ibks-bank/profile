-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id          bigserial primary key,
    created_at  timestamp     not null check ( created_at > '1970-01-01' ) default now(),
    email       text unique   not null check ( email != '' ),
    password    text          not null check ( password != '' ),
    passport_id bigint unique not null
);

create table if not exists passports
(
    id          bigserial primary key,
    series      text not null check ( length(series) = 4 ),
    number      text not null check ( length(number) = 6 ),
    first_name  text not null check ( first_name != '' ),
    middle_name text not null check ( middle_name != '' ),
    last_name   text not null check ( last_name != '' ),
    issued_by   text not null check ( issued_by != '' ),
    issued_at   date not null check ( issued_at > '1970-01-01' ),
    address     text not null check ( address != '' ),
    birthplace  text not null check ( birthplace != '' ),
    birthdate   date not null check ( birthdate > '1970-01-01' ),

    unique (series, number)
);

create table if not exists authentication_codes
(
    id         bigserial primary key,
    account_id bigint  not null references users (id),
    code       text    not null check ( code != '' ),
    expired    boolean not null default false
);

alter table users
    add constraint fk_users_passports foreign key (passport_id) references passports (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
drop table if exists passports;
-- +goose StatementEnd
