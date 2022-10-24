--File: migrations/000002_create_todo_index.up.sql
create index if not exists todos_title_idx on todos using gin(to_tsvector('simple', Title));
create index if not exists todos_description_idx on todos using gin(to_tsvector('simple', Description));