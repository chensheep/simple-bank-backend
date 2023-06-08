-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username, 
  email,
  secret_code
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails 
SET
  is_used = COALESCE(sqlc.narg('is_used'), is_used)
WHERE id = @email_id
  AND secret_code = @secret_code
  AND is_used = FALSE
  AND expired_at > now()
RETURNING *;