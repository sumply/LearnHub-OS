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

-- name: GetDomainUser :one
SELECT
    p.account_id,
    p.first_name,
    p.last_name,
    p.role,
    p.created_at
FROM account.profile AS p
WHERE p.account_id = $1;

-- name: GetDomainGroup :one
SELECT
    g.id,
    g.name,
    t.account_id AS curator_id,
    array_agg(s.account_id)::UUID[] AS student_ids
FROM school.group AS g
LEFT JOIN account.teacher AS t
    ON t.group_id = g.id
LEFT JOIN account.student AS s
    ON s.group_id = g.id
WHERE g.id = $1
GROUP BY g.id, g.name, t.account_id;

-- name: GetDomainSubject :one
SELECT 
    s.id,
    s.name
FROM school.subject AS s
WHERE s.id = $1;

-- name: ListDomainUser :many
SELECT
    p.account_id,
    p.first_name,
    p.last_name,
    p.role,
    p.created_at
FROM account.profile AS p
WHERE p.account_id IN ($1::UUID[]);

-- name: GetDomainUserByCredential :one
SELECT
    p.account_id,
    p.first_name,
    p.last_name,
    p.role,
    p.created_at
FROM account.profile AS p
JOIN account.credential AS c
    ON c.account_id = p.account_id
WHERE 
    c.email = $1 AND c.pwd_hash = $2;