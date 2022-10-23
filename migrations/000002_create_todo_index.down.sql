--File: migrations/000002_create_todo_index.down.sql
drop index if exists todos_title_idx;
drop index if exists todos_description_idx;