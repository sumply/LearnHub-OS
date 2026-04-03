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
    type
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: InsertQuizQuestionSingle :exec
INSERT INTO quiz.question_single (
    question_id,
    correct,
    options
)
VALUES (
    $1,
    $2,
    $3
);

-- name: InsertQuizQuestionMultiple :exec
INSERT INTO quiz.question_multiple (
    question_id,
    correct,
    options
)
VALUES (
    $1,
    $2,
    $3
);

-- name: InsertQuizQuestionNumeric :exec
INSERT INTO quiz.question_numeric (
    question_id,
    correct
)
VALUES (
    $1,
    $2
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