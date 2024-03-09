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

-- name: GetCustomerRootFolder :one
SELECT * FROM folder
WHERE customer_id = $1 AND parent_id IS NULL;

-- name: CreateFolderRoot :one
INSERT INTO folder (
    customer_id, title
) VALUES (
    $1, 'root'
)
RETURNING *;

-- name: CreateFolder :one
INSERT INTO folder (
    parent_id, customer_id, title
) VALUES (
    $1, $2, $3
)
RETURNING *;