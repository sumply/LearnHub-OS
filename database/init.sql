CREATE SCHEMA school;

CREATE TABLE school.subjects(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE school.groups(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
);

CREATE SCHEMA account;

CREATE TABLE account.credentials(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	pwd_hash TEXT NOT NULL
);

CREATE TABLE account.profiles(
	id UUID REFERENCES account.credentials(id) NOT NULL UNIQUE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	role INT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account.students(
	profile_id UUID REFERENCES account.profiles(id) NOT NULL UNIQUE,
	group_id UUID REFERENCES school.groups(id)
);

CREATE TABLE account.teachers(
	profile_id UUID REFERENCES account.profiles(id) NOT NULL,
	group_id UUID REFERENCES school.groups(id) NOT NULL,
	subject_id UUID REFERENCES school.subjects(id) NOT NULL,
	UNIQUE(profile_id, group_id, subject_id)
);