CREATE SCHEMA school;

-- Школьный предмет.
CREATE TABLE school.subject(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

-- Школьная группа.
CREATE TABLE school.group(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

CREATE SCHEMA account;

CREATE TYPE account.user_role AS enum('teacher', 'student', 'admin');

-- Учетные данные пользователя.
CREATE TABLE account.credential(
	account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	pwd_hash TEXT NOT NULL
);

-- Общая информация о пользователе, доступная всем.
CREATE TABLE account.profile(
	account_id UUID REFERENCES account.credential(account_id) ON DELETE CASCADE NOT NULL UNIQUE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	role account.user_role NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Расширение таблицы account.profile при роли 'student'.
CREATE TABLE account.student(
	account_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL UNIQUE,
	group_id UUID REFERENCES school.group(id) NOT NULL
);

-- Расширение таблицы account.profile при роли 'teacher'.
CREATE TABLE account.teacher(
	account_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL,
	group_id UUID REFERENCES school.group(id) NOT NULL,
	UNIQUE(account_id, group_id)
);

CREATE SCHEMA quiz;

CREATE TYPE quiz.question_type AS enum('single', 'multiple', 'numeric');

CREATE DOMAIN quiz.score AS INT DEFAULT 1 CHECK( VALUE >= 0 );

-- Общая инфомрация о квизе.
CREATE TABLE quiz.info(
	quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
	summary TEXT NOT NULL,
	subject_id UUID REFERENCES school.subject(id) NOT NULL,
	owner_id UUID REFERENCES account.profile(account_id) NOT NULL,
	max_attempts INT NOT NULL DEFAULT 1,
	total_score quiz.score NOT NULL,
	deadline TIMESTAMPTZ DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Определение доступа квиза для групп.
CREATE TABLE quiz.assignment(
	quiz_id UUID REFERENCES quiz.info(quiz_id) ON DELETE CASCADE NOT NULL,
	group_id UUID REFERENCES school.group(id) ON DELETE CASCADE NOT NULL,
	PRIMARY KEY (quiz_id, group_id) 
);

-- Хранение всех вопросов, относящихся к определенному квизу.
CREATE TABLE quiz.question(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	quiz_id UUID REFERENCES quiz.info(quiz_id) ON DELETE CASCADE NOT NULL,
	title TEXT NOT NULL,
	score quiz.score NOT NULL,
	details JSONB NOT NULL,
	position SMALLINT NOT NULL DEFAULT 0,
	UNIQUE(quiz_id, position)
);

-- Запись прохождения квиза пользователем.
CREATE TABLE quiz.attempt(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	number_attempt SMALLINT NOT NULL DEFAULT 1,
	quiz_id UUID REFERENCES quiz.info(quiz_id) ON DELETE CASCADE NOT NULL,
	user_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL,
	score quiz.score NOT NULL,
	started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ended_at TIMESTAMPTZ,
	UNIQUE(quiz_id, number_attempt, user_id)
);

-- Хранение всех ответов на вопросы в попытке.
CREATE TABLE quiz.answer(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	attempt_id UUID REFERENCES quiz.attempt(id) ON DELETE CASCADE NOT NULL,
	question_id UUID REFERENCES quiz.question(id) ON DELETE CASCADE NOT NULL,
	details JSONB NOT NULL,
	score quiz.score NOT NULL DEFAULT 0,
	is_correct BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE(attempt_id, question_id)
);

CREATE SCHEMA domain;

CREATE VIEW domain.quiz AS (
	SELECT 
		i.quiz_id,
		i.title,
		i.summary,
		i.subject_id,
		i.owner_id,
		i.max_attempts,
		i.total_score,
		i.deadline,
		i.created_at,
		array_agg(
		(
			SELECT group_id
			FROM quiz.assignment AS a
			WHERE a.quiz_id = i.quiz_id
		)
		) AS group_ids,
		COALESCE(
			json_agg(q) FILTER (WHERE q.id IS NOT NULL),
			'[]'
		) AS questions
	FROM quiz.info AS i
	LEFT JOIN quiz.question AS q
		ON i.quiz_id = q.quiz_id
	GROUP BY 
		i.quiz_id,
		i.title,
		i.summary,
		i.subject_id,
		i.owner_id,
		i.max_attempts,
		i.total_score,
		i.deadline,
		i.created_at	
);