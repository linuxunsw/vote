-- +goose Up
create table "otp_ratelimit" (
    "zid" text primary key,
    "count" integer not null default 0,
    "window_start" timestamp (3) with time zone not null
);

-- +goose Down
drop table "otp_ratelimit";