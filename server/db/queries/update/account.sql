
-- name: UpdateAccountProfile :exec
UPDATE account.profile
SET 
    first_name = $1,
    last_name = $2,
    role = $3
WHERE account_id = $4;
