CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firstName VARCHAR(40) NOT NULL,
    secondName VARCHAR(40) NOT NULL,
    lastName VARCHAR(40),
    email TEXT UNIQUE NOT NULL,
    passwordHash TEXT NOT NULL,
    role SMALLINT NOT NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id),
    is_visible BOOLEAN DEFAULT TRUE,
    requires_submission BOOLEAN DEFAULT FALSE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task
