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

CREATE TABLE account.credential(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	pwd_hash TEXT NOT NULL
);

CREATE TABLE account.profile(
	id UUID REFERENCES account.credential(id) ON DELETE CASCADE NOT NULL UNIQUE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	role INT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account.student(
	profile_id UUID REFERENCES account.profile(id) ON DELETE CASCADE NOT NULL UNIQUE,
	group_id UUID REFERENCES school.group(id)
);

CREATE TABLE account.teacher(
	profile_id UUID REFERENCES account.profile(id) ON DELETE CASCADE NOT NULL,
	group_id UUID REFERENCES school.group(id) NOT NULL,
	subject_id UUID REFERENCES school.subject(id) NOT NULL,
	UNIQUE(profile_id, group_id, subject_id)
);

CREATE SCHEMA quiz;

CREATE TABLE quiz.info(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
	summary TEXT NOT NULL,
	subject_id UUID REFERENCES school.subject(id) NOT NULL,
	owner_id UUID REFERENCES account.profile(id) NOT NULL,
	total_score SMALLINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE quiz.content(
	quiz_id UUID REFERENCES quiz.info(id) ON DELETE CASCADE NOT NULL UNIQUE,
	content JSONB NOT NULL
);

CREATE TABLE quiz.progress(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	quiz_id UUID REFERENCES quiz.info(id) ON DELETE CASCADE NOT NULL,
	user_id UUID REFERENCES account.profile(id) ON DELETE CASCADE NOT NULL,
	content JSONB NOT NULL,
	score SMALLINT NOT NULL DEFAULT 0,
	started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ended_at TIMESTAMPTZ,
	UNIQUE(quiz_id, user_id)
);