
-- name: DeleteSchoolGroup :exec
DELETE FROM school.group
WHERE id = $1;

-- name: DeleteSchoolSubject :exec
DELETE FROM school.subject
WHERE id = $1;
