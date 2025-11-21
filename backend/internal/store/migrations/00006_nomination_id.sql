-- +goose Up
alter table nominations add column nomination_id uuid;

alter table nominations drop constraint if exists nominations_pkey;

alter table nominations
add constraint nominations_election_candidate_key unique (
    election_id, candidate_zid
);

alter table nominations
add constraint nominations_pkey primary key (nomination_id);

-- +goose Down
alter table nominations
drop constraint if exists nominations_pkey;

alter table nominations
drop constraint if exists nominations_election_candidate_key;

alter table nominations
drop column if exists nomination_id;

alter table nominations
add constraint nominations_pkey primary key (election_id, candidate_zid);
