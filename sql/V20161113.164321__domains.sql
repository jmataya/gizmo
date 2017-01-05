create domain generic_string text check (length(value) <= 255);
create domain generic_timestamp timestamp without time zone default (now() at time zone 'utc');
create domain generic_timestamp_null timestamp without time zone null;
