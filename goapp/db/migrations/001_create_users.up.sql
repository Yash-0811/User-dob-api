-- db/migrations/001_create_users.up.sql

CREATE TABLE IF NOT EXISTS users (
    id   SERIAL      PRIMARY KEY,
    name TEXT        NOT NULL,
    dob  DATE        NOT NULL
);
