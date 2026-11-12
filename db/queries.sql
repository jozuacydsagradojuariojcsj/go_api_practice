-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at;

-- name: LockUsers :many
SELECT id, balance FROM users
WHERE id = ANY(@ids::uuid[])
ORDER BY id
FOR UPDATE;

-- name: AddBalance :one
UPDATE users
SET balance = balance + @amount
WHERE id = @id
RETURNING *;

-- name: RecordTransaction :one
INSERT INTO transactions (from_user_id, to_user_id, amount)
VALUES($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET name = @name, email = @email
WHERE id = @id
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = @id;

-- name: ListUserTransaction :many
SELECT * FROM transactions
WHERE from_user_id =sqlc.arg(user_id)::uuid OR to_user_id = sqlc.arg(user_id)::uuid
ORDER BY created_at DESC
LIMIT 50;