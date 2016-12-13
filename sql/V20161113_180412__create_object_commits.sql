create table object_commits (
  id serial primary key,
  form_id integer not null references object_forms(id) on update restrict on delete restrict,
  shadow_id integer not null references object_shadows(id) on update restrict on delete restrict,
  previous_id integer null references object_commits(id) on update restrict on delete restrict,

  created_at generic_timestamp,

  foreign key (form_id) references object_forms(id) on update restrict on delete restrict,
  foreign key (shadow_id) references object_shadows(id) on update restrict on delete restrict
);

create index object_commits_object_form_idx on object_commits (form_id);
create index object_commits_object_shadow_idx on object_commits (shadow_id);
