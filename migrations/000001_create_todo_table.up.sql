--File: migrations/000001_create_todo_table.up.sql
CREATE TABLE IF NOT EXISTS todos(
    ID bigserial PRIMARY KEY,
    CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    Title text not null,
    Description text not null,
    Done text not null,
    Version integer not null DEFAULT 1
);