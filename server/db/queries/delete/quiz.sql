-- name: DeleteQuizInfo :exec
DELETE FROM quiz.info
WHERE quiz_id = $1;
