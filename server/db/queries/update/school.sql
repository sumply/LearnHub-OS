-- name: UpdateSchoolGroup :exec
UPDATE school.group
SET 
    name = $1
WHERE id = $2;

-- name: UpdateSchoolSubject :exec
UPDATE school.subject
SET 
    name = $1
WHERE id = $2;