-- name: DeleteAccountCredential :exec
DELETE FROM account.credential
WHERE account_id = $1;
