
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
WHERE p.account_id = ANY($1::UUID[]);

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

-- name: GetDomainQuiz :one
SELECT 
    i.quiz_id,
    i.title,
    i.summary,
    i.owner_id,
    i.subject_id,
    i.max_attempts,
    i.deadline,
    i.total_score,
    i.created_at,
    COALESCE(q.questions, '[]')::JSONB AS questions,
    COALESCE(a.group_ids, '{}')::UUID[] AS group_ids
FROM quiz.info AS i
LEFT JOIN LATERAL (
    SELECT json_agg(
        json_build_object(
            'id', q.id,
            'title', q.title,
            'score', q.score,
            'type', q.type,
            'single_correct', COALESCE(s.correct, '')::TEXT,
            'single_options', COALESCE(s.options, '{}')::TEXT[],
            'multiple_correct', COALESCE(m.correct, '{}')::TEXT[],
            'multiple_options', COALESCE(m.options, '{}')::TEXT[],
            'numeric_correct', COALESCE(n.correct, 0)::FLOAT
        )
    ) AS questions
    FROM quiz.question AS q
    LEFT JOIN quiz.question_single AS s ON s.question_id = q.id
    LEFT JOIN quiz.question_multiple AS m ON m.question_id = q.id
    LEFT JOIN quiz.question_numeric AS n ON n.question_id = q.id
    WHERE q.quiz_id = i.quiz_id
) AS q ON true
LEFT JOIN LATERAL (
    SELECT array_agg(a.group_id)::UUID[] AS group_ids
    FROM quiz.assignment AS a
    WHERE a.quiz_id = i.quiz_id 
) AS a ON true
WHERE i.quiz_id = $1;

-- name: GetDomainQuizzes :many
SELECT 
    i.quiz_id,
    i.title,
    i.summary,
    i.owner_id,
    i.subject_id,
    i.max_attempts,
    i.deadline,
    i.total_score,
    i.created_at,
    COALESCE(a.group_ids, '{}')::UUID[] AS group_ids
FROM quiz.info AS i
LEFT JOIN LATERAL (
    SELECT array_agg(a.group_id)::UUID[] AS group_ids
    FROM quiz.assignment AS a
    WHERE a.quiz_id = i.quiz_id 
) AS a ON true;

-- name: GetDomainAttempts :many
SELECT 
    a.id,
    a.number_attempt,
    a.quiz_id,
    a.user_id,
    a.score,
    a.started_at,
    a.ended_at
FROM quiz.attempt AS a;

-- name: GetDomainAttempt :one
SELECT 
    att.id,
    att.number_attempt,
    att.quiz_id,
    att.user_id,
    att.score,
    att.started_at,
    att.ended_at,
    COALESCE(a.answers, '[]')::JSONB AS answers
FROM quiz.attempt AS att
LEFT JOIN LATERAL (
    SELECT json_agg(
        json_build_object(
            'id', a.id,
            'question_id', a.question_id,
            'score', a.score,
            'is_correct', a.is_correct,
            'type', q.type,
            'single_selected_answer', COALESCE(s.selected_answer, '')::TEXT,
            'multiple_selected_answer', COALESCE(m.selected_answer, '{}')::TEXT[],
            'numeric_selected_answer', COALESCE(n.selected_answer, 0)::FLOAT
        )
    ) AS answers
    FROM quiz.answer AS a
    INNER JOIN quiz.question AS q ON q.id = a.question_id
    LEFT JOIN quiz.answer_single AS s ON s.answer_id = a.id
    LEFT JOIN quiz.answer_multiple AS m ON m.answer_id = a.id
    LEFT JOIN quiz.answer_numeric AS n ON n.answer_id = a.id
    WHERE a.attempt_id = att.id
) as a ON true
WHERE att.id = $1;

