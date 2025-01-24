-- +goose Up
CREATE TABLE users (
  id serial PRIMARY KEY,
  name varchar(128) NOT NULL,
  email varchar(128) NOT NULL,
  password text NOT NULL,
  created_at timestamp NOT NULL default now(),
  updated_at timestamp
);

-- +goose Down
DROP TABLE users;