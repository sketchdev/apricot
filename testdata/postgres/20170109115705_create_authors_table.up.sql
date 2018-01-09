create table authors (
  id SERIAL PRIMARY KEY NOT NULL,
  name varchar(50),
  created_at timestamp default now()
);