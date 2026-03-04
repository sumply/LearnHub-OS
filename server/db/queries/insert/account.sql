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

-- name: InsertAccountStudent :exec
INSERT INTO account.student (
    account_id,
    group_id
)
VALUES (
    $1,
    $2
);

-- name: InsertAccountTeacher :exec
INSERT INTO account.teacher (
    account_id,
    group_id
)
VALUES (
    $1,
    $2
);