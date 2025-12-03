CREATE EXTENSION citext;

CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS school;
CREATE SCHEMA IF NOT EXISTS quiz;

CREATE TYPE users.role_enum AS enum('student', 'teacher', 'admin', 'root');
CREATE TYPE quiz.question_enum AS enum('one', 'multiple', 'text');

-- Секретная информация, используемая при входе в аккаунт.
CREATE TABLE IF NOT EXISTS users.auth_data (
  login TEXT PRIMARY KEY,
  email CITEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL
);

-- Общая информация о пользователе.
CREATE TABLE IF NOT EXISTS users.profile (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(40) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    middle_name VARCHAR(40),
    role users.role_enum NOT NULL,
	login_ref TEXT REFERENCES users.auth_data(login) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Классы школы. Например, 11А
CREATE TABLE IF NOT EXISTS school.groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Школьные предметы
CREATE TABLE IF NOT EXISTS school.subjects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Пользователь относится к классу.
CREATE TABLE IF NOT EXISTS users.groups (
    user_id INTEGER REFERENCES users.profile(id) ON DELETE CASCADE,
    group_id INTEGER REFERENCES school.groups(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, group_id)
);

-- Общая информация о заданиях.
CREATE TABLE IF NOT EXISTS quiz.quizzes (
  id SERIAL PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  summary VARCHAR(230),
  subject_id INTEGER REFERENCES school.subjects(id) ON DELETE CASCADE,
  teacher_id INTEGER REFERENCES users.profile(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Группы, которым доступен квизис
CREATE TABLE IF NOT EXISTS quiz.groups (
	quizz_id INTEGER REFERENCES quiz.quizzes(id) ON DELETE CASCADE,
	group_id INTEGER REFERENCES school.groups(id) ON DELETE CASCADE,
	UNIQUE(quizz_id, group_id)
);

-- Вопросы к заданиям.
CREATE TABLE IF NOT EXISTS quiz.questions (
	id SERIAL PRIMARY KEY,
	quizz_id INTEGER REFERENCES quiz.quizzes(id) ON DELETE CASCADE,
	question VARCHAR(230) NOT NULL
);

-- Ответы на вопросы к заданиям.
CREATE TABLE IF NOT EXISTS quiz.answer_options (
	id SERIAL PRIMARY KEY,
	question_id INTEGER REFERENCES quiz.questions(id) ON DELETE CASCADE,
	option_text VARCHAR(230) NOT NULL,
	is_correct BOOLEAN NOT NULL DEFAULT FALSE
);

-- Ответы пользователя на вопросы заданий.
CREATE TABLE IF NOT EXISTS quiz.selected_answers (
	question_id INTEGER REFERENCES quiz.questions(id) ON DELETE CASCADE,
	user_id INTEGER REFERENCES users.profile(id) ON DELETE CASCADE,
	answer_id INTEGER REFERENCES quiz.answer_options(id) ON DELETE CASCADE,
	UNIQUE(question_id, answer_id)
);

-- Общая информация о пользователе, которому доступно задание.
CREATE TABLE IF NOT EXISTS quiz.progress (
  user_id INTEGER REFERENCES users.profile(id) ON DELETE CASCADE,
  quizz_id INTEGER REFERENCES quiz.quizzes(id) ON DELETE CASCADE,
  completed BOOL NOT NULL DEFAULT false,
  completion_time TIMESTAMP NOT NULL,
  UNIQUE(user_id, quizz_id)
);

/*
-- Представление количества вопросов в квизе.
CREATE OR REPLACE VIEW quiz.questions_count AS (
	SELECT 
		quizzes.id, 
		quizzes.name, 
		COUNT(*) AS question_count
	FROM 
		quizzes
	LEFT JOIN 
		quizz_questions AS questions
		ON questions.quizz_id = quizzes.id
	GROUP BY quizzes.id, quizzes.name
);

-- Представление количества верно решенных вопросов из квиза.
CREATE OR REPLACE VIEW quiz.user_correct_answers_count AS (
	SELECT 
		u.id, 
		u.first_name, 
		u.last_name, 
		u.middle_name,
		COUNT(DISTINCT a.question_id) AS correct_answer_count
	FROM user_public AS u
	JOIN
		user_quizz_answer AS a
		ON a.user_id = u.id
	JOIN 
		quizz_answer_options AS o
		ON o.id = a.answer_id
	WHERE
		o.is_correct = true
	GROUP BY 
		u.id
);
*/