--File: migrations/000002_create_todo_index.up.sql
create index if not exists todos_title_idx on todo_list using gin(to_tsvector('simple', title));
create index if not exists todos_description_idx on todo_list using gin(to_tsvector('simple', description));