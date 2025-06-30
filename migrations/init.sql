-- Create user 'postgres' with password
DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres'
   ) THEN
CREATE ROLE postgres WITH LOGIN PASSWORD 'postgres';
ALTER ROLE postgres CREATEDB;
END IF;
END
$$;

-- Create database 'shortener' assigned to user 'postgres'
DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_database WHERE datname = 'shortener'
   ) THEN
      CREATE DATABASE shortener OWNER postgres;
END IF;
END
$$;

-- Connection to database 'shortener' and create table
\connect shortener

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_id VARCHAR(16) UNIQUE NOT NULL,
    long_url TEXT NOT NULL
);

-- Add accesses
GRANT ALL PRIVILEGES ON TABLE urls TO postgres;
GRANT ALL PRIVILEGES ON DATABASE shortener TO postgres;
