-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS social;

CREATE TABLE IF NOT EXISTS social.users (
    id varchar PRIMARY KEY NOT NULL ,
    first_name varchar NOT NULL ,
    second_name varchar NOT NULL ,
    age integer NOT NULL,
    birthdate date NOT NULL ,
    biography varchar,
    city varchar,
    password varchar NOT NULL

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS social CASCADE;
-- +goose StatementEnd
