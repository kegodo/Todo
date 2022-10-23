--File: migrations/000001_create_todo_table.up.sql
create table if not exists todo_list(
    id int not null,
    created_at timestamp(0) with time zone not null default now(),
    title text not null,
    description text not null,
    done boolean not null,
    version int not null default 1
);