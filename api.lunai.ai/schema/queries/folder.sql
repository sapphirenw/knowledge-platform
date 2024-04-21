-- name: GetFolder :one
SELECT * FROM folder
WHERE id = $1 LIMIT 1;

-- name: GetFolderWithName :one
SELECT * FROM folder
WHERE customer_id = $1 AND title = $2
LIMIT 1;

-- name: GetFoldersByCustomer :many
SELECT * FROM folder
WHERE customer_id = $1;

-- name: GetFoldersFromParent :many
SELECT * FROM folder
WHERE parent_id = $1;

-- name: GetRootFoldersByCustomer :many
SELECT * FROM folder
WHERE customer_id = $1 AND parent_id IS NULL;

-- name: CreateFolderRoot :one
INSERT INTO folder (
    customer_id, title
) VALUES (
    $1, 'root'
)
ON CONFLICT (customer_id, COALESCE(parent_id, -1), title) DO UPDATE
SET updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: CreateFolder :one
INSERT INTO folder (
    parent_id, customer_id, title
) VALUES (
    $1, $2, $3
)
ON CONFLICT (customer_id, COALESCE(parent_id, -1), title) DO UPDATE
SET updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeleteFoldersOlderThan :exec
DELETE FROM folder
WHERE customer_id = $1
AND updated_at < $2;

-- name: GetFoldersOlderThan :many
SELECT * FROM folder
WHERE customer_id = $1
AND updated_at < $2;