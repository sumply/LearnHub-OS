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