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
