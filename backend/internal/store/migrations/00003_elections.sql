-- +goose Up
create table "election_member_list" (
    "zid" text,
    "election_id" uuid,

    primary key ("election_id", "zid")
);

-- NO_ELECTION is external
create type "election_status" as enum (
    'CLOSED',
    'NOMINATION_OPEN',
    'NOMINATION_CLOSED',
    'VOTING_OPEN',
    'VOTING_CLOSED',
    'RESULTS',
    'END'
);

create table "elections" (
    "election_id" uuid primary key,
    "name" text not null,
    
    "state" election_status not null,

    "created_at" timestamptz not null,

    "nominations_open_at" timestamptz,
    "nominations_close_at" timestamptz,

    "voting_open_at" timestamptz,
    "voting_close_at" timestamptz,

    "results_published_at" timestamptz,
    "ended_at" timestamptz
);

create index "elections_created_at_idx" on "elections" ("created_at");

-- +goose Down
drop table "election_member_list";
drop table "elections";
drop type "election_status";