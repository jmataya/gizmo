create table object_forms (
  id serial primary key,
  kind generic_string not null,
  attributes jsonb,

  created_at generic_timestamp,
  updated_at generic_timestamp
);

create index object_forms_kind_idx on object_forms (kind);
  
