-- name: GetFolder :one
SELECT * FROM folder
WHERE id = $1 LIMIT 1;

-- name: GetFoldersByCustomer :many
SELECT * FROM folder
WHERE customer_id = $1;

-- name: GetFoldersFromParent :many
SELECT * FROM folder
WHERE parent_id = $1;

-- name: GetCustomerRootFolder :one
SELECT * FROM folder
WHERE customer_id = $1 AND parent_id = NULL;

-- name: GetDocument :one
SELECT * FROM document
WHERE id = $1 LIMIT 1;

-- name: GetDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1;

-- name: GetDocumentsFromParent :many
SELECT * FROM document
where parent_id = $1;