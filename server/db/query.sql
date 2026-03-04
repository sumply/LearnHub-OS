-- ========================================
-- SCHOOL SCHEMA
-- ========================================

-- name: InsertSchoolSubject :exec
INSERT INTO school.subject (
    id,
    name
)
VALUES (
    $1, 
    $2
);

-- name: InsertSchoolGroup :exec
INSERT INTO school.group (
    id,
    name
)
VALUES (
    $1, 
    $2
);

-- ========================================
-- ACCOUNT SCHEMA
-- ========================================

-- name: InsertAccountCredential :exec
INSERT INTO account.credential (
    account_id,
    email,
    pwd_hash
)
VALUES (
    $1,
    $2,
    $3
);

-- name: InsertAccountProfile :exec
INSERT INTO account.profile (
    account_id,
    first_name,
    last_name,
    role,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: InsertStudent :exec
INSERT INTO account.student (
    account_id,
    group_id
)
VALUES (
    $1,
    $2
);

-- name: InsertTeacher :exec
INSERT INTO account.teacher (
    account_id,
    group_id
)
VALUES (
    $1,
    $2
);

-- ========================================
-- QUIZ SCHEMA
-- ========================================

-- name: InsertQuizInfo :exec
INSERT INTO quiz.info (
    quiz_id,
    title,
    summary,
    subject_id,
    owner_id,
    max_attempts,
    total_score,
    deadline,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);

-- name: InsertQuizQuestion :exec
INSERT INTO quiz.question (
    id,
    quiz_id,
    title,
    score,
    details
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: InsertQuizAttempt :exec
INSERT INTO quiz.attempt (
    id,
    quiz_id,
    user_id,
    started_at
)
VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: InsertQuizAnswer :exec
INSERT INTO quiz.answer (
    id,
    attempt_id,
    question_id,
    details
)
VALUES (
    $1,
    $2,
    $3,
    $4
);


-- name: InsertQuizAssigment :exec
INSERT INTO quiz.assignment(
    quiz_id,
    group_id
)
VALUES (
    $1,
    $2
);

-- name: GetDetailedUsers :many
SELECT 
	p.account_id, 
	p.first_name, 
	p.last_name, 
	p.role, 
	p.created_at,
	s.group_id AS s_group_id,
	array_remove(array_agg(t.group_id), NULL)::UUID[] AS t_group_ids
FROM account.profile AS p
LEFT JOIN account.student AS s
	ON s.account_id = p.account_id
LEFT JOIN account.teacher AS t
	ON t.account_id = p.account_id
GROUP BY 
	p.account_id, 
	p.first_name,
	p.last_name,
	p.role,
	p.created_at,
	s.group_id;

-- name: GetQuiz :one
SELECT 
	i.quiz_id AS quiz_id, 
	i.title AS quiz_title,
	i.summary AS quiz_summary,
	i.total_score::INT AS quiz_total_score,
	i.deadline AS quiz_deadline,
	i.max_attempts AS quiz_max_attempts,
	i.created_at AS quiz_created_at,
    json_agg(q.*)::JSONB AS quiz_questions,
	p.account_id AS owner_id,
	p.first_name AS owner_first_name,
	p.last_name AS owner_last_name,
	p.role AS owner_role,
	p.created_at AS owner_created_at,
	s.id AS subject_id,
	s.name AS subject_name
FROM quiz.info AS i
JOIN quiz.question AS q
    ON q.quiz_id = i.quiz_id
JOIN account.profile AS p 
	ON p.account_id = i.owner_id
JOIN school.subject AS s
	ON s.id = i.subject_id
WHERE i.quiz_id = $1
GROUP BY
	i.quiz_id, 
	i.title,
	i.summary,
	i.total_score,
	i.deadline,
	i.max_attempts,
	i.created_at,
	p.account_id,
	p.first_name,
	p.last_name ,
	p.role,
	p.created_at,
	s.id,
	s.name;

-- name: GetQuizItem :many
SELECT 
	i.quiz_id AS quiz_id, 
	i.title AS quiz_title,
	i.summary AS quiz_summary,
	i.total_score::INT AS quiz_total_score,
	i.created_at AS quiz_created_at,
    i.deadline AS quiz_deadline,
    i.max_attempts::INT AS quiz_max_attempts,
	p.account_id AS owner_id,
	p.first_name AS owner_first_name,
	p.last_name AS owner_last_name,
	p.role AS owner_role,
	p.created_at AS owner_created_at,
	s.id AS subject_id,
	s.name AS subject_name
FROM quiz.info AS i
JOIN account.profile AS p
	ON p.account_id = i.owner_id
JOIN school.subject AS s
	ON s.id = i.subject_id;

-- name: GetDomainQuiz :one
SELECT 
	i.quiz_id AS quiz_id,
	i.title AS quiz_title,
	i.summary AS quiz_summary,
	i.owner_id AS quiz_owner_id,
	i.subject_id AS quiz_subject_id,
	i.total_score::INT AS quiz_total_score,
	i.deadline AS quiz_deadline,
	i.max_attempts::INT AS quiz_max_attempts,
	i.created_at AS quiz_created_at,
	json_agg(q.*)::JSONB AS quiz_questions
FROM quiz.info AS i
JOIN quiz.question AS q
	ON q.quiz_id = i.quiz_id
WHERE i.quiz_id = $1
GROUP BY 
    i.quiz_id,
    i.title,
    i.summary,
    i.owner_id,
    i.subject_id,
    i.total_score,
    i.deadline,
    i.max_attempts,
    i.created_at;

-- name: GetDomainAttempt :one
SELECT 
    attempt.id AS attempt_id,
    attempt.quiz_id AS attempt_quiz_id,
    attempt.user_id AS attempt_user_id,
    attempt.score::INT AS attempt_score_id,
    attempt.started_at AS attempt_started_at,
    attempt.ended_at AS attempt_endend_at,
    json_agg(answer.*) AS attempt_answers
FROM quiz.attempt AS attempt
JOIN quiz.answer AS answer
    ON answer.attempt_id = attempt.id
GROUP BY 
	attempt.id,
	attempt.quiz_id,
    attempt.user_id,
    attempt.score,
    attempt.started_at,
    attempt.ended_at;

-- name: GetFinishedAttempt :one
SELECT 
    attempt.id AS attempt_id,
    attempt.score::INT AS attempt_score,
    attempt.started_at AS attempt_started_at,
    attempt.ended_at AS attempt_ended_at,
    (
		SELECT json_agg(answer.*) 
		FROM quiz.answer 
		WHERE attempt_id = $1
	) AS attempt_answers,

    a_user.account_id AS user_id,
    a_user.first_name AS user_first_name,
    a_user.last_name AS user_last_name,
    a_user.role AS user_role,
    a_user.created_at AS user_created_at,

    q_info.quiz_id AS quiz_id,
    q_info.title AS quiz_title,
    q_info.summary AS quiz_summary,
    q_info.total_score::INT AS quiz_total_score,
    q_info.deadline AS quiz_deadline,
    q_info.max_attempts AS quiz_max_attempts,
    q_info.created_at AS quiz_created_at,
    (
		SELECT json_agg(question.*)
		FROM quiz.question
		WHERE quiz_id = q_info.quiz_id
	)AS quiz_questions,

    q_owner.account_id AS owner_id,
    q_owner.first_name AS owner_first_name,
    q_owner.last_name AS owner_last_name,
    q_owner.role AS owner_role,
    q_owner.created_at AS owner_created_at,

    q_subject.id AS subject_id,
    q_subject.name AS subject_name

FROM quiz.attempt AS attempt

JOIN account.profile AS a_user 
    ON a_user.account_id = attempt.user_id

JOIN quiz.info AS q_info
    ON q_info.quiz_id = attempt.quiz_id

JOIN account.profile AS q_owner 
    ON q_owner.account_id = q_info.owner_id
    
JOIN school.subject AS q_subject
    ON q_subject.id = q_info.subject_id

WHERE attempt.id = $1;

-- name: GetDomainTeacher :one 
SELECT
    p.account_id,
    p.first_name,
    p.last_name,
    p.role,
    p.created_at,
    array_agg(t.group_id)::UUID[] AS group_ids
FROM account.profile AS p
LEFT JOIN account.teacher AS t
    ON t.account_id = p.account_id
WHERE p.account_id = $1
GROUP BY
    p.account_id,
    p.first_name,
    p.last_name,
    p.role,
    p.created_at;

-- name: UpdateQuizAttempt :exec
UPDATE quiz.attempt 
SET 
    score = $1,
    ended_at = $2
WHERE id = $3;

-- name: UpdateQuizAnswer :exec
UPDATE quiz.answer 
SET 
    details = $1,
    score = $2,
    is_correct = $3
WHERE id = $4;