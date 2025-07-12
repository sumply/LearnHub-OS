drop database if exists main;
CREATE DATABASE main;
\c main
-- use database main;
drop TABLE if exists users CASCADE;
create type role as enum('student', 'teacher', 'admin', 'root');
create table
users (
    id SERIAL PRIMARY KEY,
    firstName varchar(40) not null,
    secondtName varchar(40) not null,
    lastName varchar(40) null,
    email TEXT UNIQUE NOT NULL,
    passwordHash TEXT NOT NULL,
    role role DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

drop TABLE if exists categories CASCADE;
create table
categories (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT
);

drop TABLE if exists tasks CASCADE;
create table
tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id),
    is_visible BOOLEAN DEFAULT TRUE,
    requires_submission BOOLEAN DEFAULT FALSE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

drop TABLE if exists task_materials CASCADE;
create table
task_materials (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    material_type TEXT CHECK (material_type IN ('video', 'presentation', 'test', 'file')),
    content TEXT, -- ссылка или текстовое поле, или имя файла
    file_name TEXT, -- оригинальное имя файла
    mime_type TEXT
);

drop TABLE if exists task_submissions CASCADE;
create table
task_submissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    task_id INTEGER REFERENCES tasks(id),
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT CHECK (status IN ('submitted', 'marked')),
    mark INTEGER CHECK (mark BETWEEN 0 AND 5),
    comment TEXT,
    reviewed_by INTEGER REFERENCES users(id)
);

drop TABLE if exists submission_files CASCADE;
create table
submission_files (
    id SERIAL PRIMARY KEY,
    submission_id INTEGER REFERENCES task_submissions(id) ON DELETE CASCADE,
    file_name TEXT,
    file_path TEXT,
    mime_type TEXT
);

drop TABLE if exists groups CASCADE;
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

drop TABLE if exists user_groups CASCADE;
CREATE TABLE user_groups (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, group_id)
);