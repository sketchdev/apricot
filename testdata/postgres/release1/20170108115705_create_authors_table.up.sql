create table authors (
  id serial primary key not null,
  name varchar(50) not null ,
  created_at timestamp default now()
);