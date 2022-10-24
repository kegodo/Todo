--File: migrations/000001_create_todo_table.up.sql
CREATE TABLE IF NOT EXISTS todos(
    ID int not null,
    CreatedAt timestamp(0) with time zone not null default now(),
    Title text not null,
    Description text not null,
    Done boolean not null,
    Version int not null default 1
);