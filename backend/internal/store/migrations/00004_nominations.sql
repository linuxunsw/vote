-- +goose Up
create table "nominations" (
    "candidate_zid" text,
    "election_id" uuid,
    
    "candidate_name" text not null,
    "contact_email" text not null,
    "discord_username" text,

    "executive_roles" text[] not null,
    "url" text,
    "candidate_statement" text not null,

    "created_at" timestamptz not null,
    "updated_at" timestamptz not null,

    primary key ("election_id", "candidate_zid")
);

-- +goose Down
drop table "nominations";