CREATE TABLE IF NOT EXISTS documentations (
    id serial PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    content bytea NOT NULL
);
