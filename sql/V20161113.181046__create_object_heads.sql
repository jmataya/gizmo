create table object_heads (
  id serial primary key,
  context_id integer not null references object_contexts(id) on update restrict on delete restrict,
  commit_id integer not null references object_commits(id) on update restrict on delete restrict,

  created_at generic_timestamp,
  updated_at generic_timestamp,
  archived_at generic_timestamp,

  foreign key (context_id) references object_contexts(id) on update restrict on delete restrict,
  foreign key (commit_id) references object_commits(id) on update restrict on delete restrict
);

create index object_heads_object_commit_idx on object_heads (context_id);
