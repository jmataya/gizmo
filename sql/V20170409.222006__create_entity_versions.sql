create table entity_versions (
  id serial primary key,
  parent_id integer null references entity_versions(id) on update restrict on delete restrict,
  content_commit_id integer not null references object_commits(id) on update restrict on delete restrict,

  created_at generic_timestamp,

  foreign key (content_commit_id) references object_commits(id) on update restrict on delete restrict
);
