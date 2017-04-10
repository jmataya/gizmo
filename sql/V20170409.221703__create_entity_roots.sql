create table entity_roots (
  id serial primary key,
  kind generic_string not null,

  created_at generic_timestamp,
  archived_at generic_timestamp_null
);

create index entity_roots_kind_idx on entity_roots (kind)
