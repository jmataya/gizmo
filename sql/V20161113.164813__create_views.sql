create table views (
  id serial primary key,
  name generic_string not null,
  attributes jsonb,

  created_at generic_timestamp,
  updated_at generic_timestamp
);
