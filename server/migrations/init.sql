CREATE SCHEMA school;

CREATE TABLE school.subject(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE school.group(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

CREATE SCHEMA account;

CREATE TYPE account.user_role AS enum('teacher', 'student', 'admin');

CREATE TABLE account.credential(
	account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	pwd_hash TEXT NOT NULL
);

CREATE TABLE account.profile(
	account_id UUID REFERENCES account.credential(account_id) ON DELETE CASCADE NOT NULL UNIQUE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	role account.user_role NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account.student(
	account_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL UNIQUE,
	group_id UUID REFERENCES school.group(id) NOT NULL
);

CREATE TABLE account.teacher(
	account_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL,
	group_id UUID REFERENCES school.group(id) NOT NULL,
	subject_id UUID REFERENCES school.subject(id) NOT NULL,
	UNIQUE(account_id, group_id, subject_id)
);

CREATE SCHEMA quiz;

CREATE TYPE quiz.question_type AS enum('single', 'multiple', 'numeric');

CREATE DOMAIN quiz.score AS INT DEFAULT 1 CHECK( VALUE >= 0 );

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

CREATE TABLE quiz.question(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	quiz_id UUID REFERENCES quiz.info(quiz_id) ON DELETE CASCADE NOT NULL,
	title TEXT NOT NULL,
	score quiz.score NOT NULL,
	details JSONB NOT NULL
);

CREATE TABLE quiz.attempt(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	quiz_id UUID REFERENCES quiz.info(quiz_id) ON DELETE CASCADE NOT NULL,
	user_id UUID REFERENCES account.profile(account_id) ON DELETE CASCADE NOT NULL,
	score quiz.score NOT NULL,
	started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ended_at TIMESTAMPTZ,
	UNIQUE(quiz_id, user_id)
);

CREATE TABLE quiz.answer(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	attempt_id UUID REFERENCES quiz.attempt(id) ON DELETE CASCADE NOT NULL,
	question_id UUID REFERENCES quiz.question(id) ON DELETE CASCADE NOT NULL UNIQUE,
	details JSONB NOT NULL,
	score quiz.score NOT NULL,
	is_correct BOOLEAN NOT NULL
);