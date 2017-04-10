create table entity_heads (
  id serial primary key,
  root_id integer not null references entity_roots(id) on update restrict on delete restrict,
  context_id integer not null references object_contexts(id) on update restrict on delete restrict,
  version_id integer not null references entity_versions(id) on update restrict on delete restrict,

  created_at generic_timestamp,
  updated_at generic_timestamp,
  archived_at generic_timestamp_null,

  foreign key (root_id) references entity_roots(id) on update restrict on delete restrict,
  foreign key (context_id) references object_contexts(id) on update restrict on delete restrict,
  foreign key (version_id) references entity_versions(id) on update restrict on delete restrict
);

create index entity_heads_object_context_idx on entity_heads (context_id);
