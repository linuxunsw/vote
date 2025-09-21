-- +goose Up
create table "ballots" (
    "election_id" uuid,
    "zid" text,
    
	-- {"position": "candidate_zid", "position2": "candidate_zid2", ...}
    "positions" JSONB not null,

    "created_at" timestamptz not null,
    "updated_at" timestamptz not null,

    primary key ("election_id", "zid")
);

-- +goose Down
drop table "ballots";