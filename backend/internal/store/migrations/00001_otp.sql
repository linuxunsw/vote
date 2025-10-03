-- +goose Up
create table otp (
    zid text primary key,
    code_hash text not null,
    retry_amount integer not null default 0,
    created_at timestamp (3) with time zone not null
);

-- +goose Down
drop table otp;
