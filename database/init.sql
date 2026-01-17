CREATE EXTENSION citext;

CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.role AS ENUM('none', 'teacher', 'student');
CREATE TYPE users.access AS ENUM('user', 'admin');

CREATE DOMAIN users.email AS TEXT CHECK (VALUE ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$');
-- Секретная информация, используемая при входе в аккаунт.
CREATE TABLE IF NOT EXISTS users.credentials (
	id SERIAL PRIMARY KEY,
	login TEXT NOT NULL UNIQUE,
	email users.email UNIQUE NOT NULL,
	password_hash TEXT NOT NULL
);

CREATE DOMAIN users.name AS VARCHAR(40) CHECK (LENGTH(VALUE) > 2);
-- Общая информация о пользователе.
CREATE TABLE IF NOT EXISTS users.profiles (
    id BIGINT REFERENCES users.credentials(id) NOT NULL UNIQUE,
    first_name users.name NOT NULL,
    last_name users.name NOT NULL,
    middle_name VARCHAR(40) NOT NULL,
    role users.role NOT NULL,
	access users.access NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);