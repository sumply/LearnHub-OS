-- name: FinishAttempt :exec
UPDATE quiz.attempt
SET 
    score = $1,
    ended_at = $2
WHERE id = $3;

-- name: FinishAnswer :exec
UPDATE quiz.answer
SET
    score = $1,
    is_correct = $2
WHERE id = $3;

-- name: FinishSingleAnswer :exec
UPDATE quiz.answer_single
SET
    selected_answer = $1
WHERE answer_id = $2;

-- name: FinishMultipleAnswer :exec
UPDATE quiz.answer_multiple
SET
    selected_answer = $1
WHERE answer_id = $2;

-- name: FinishNumericAnswer :exec
UPDATE quiz.answer_numeric
SET
    selected_answer = $1
WHERE answer_id = $2;