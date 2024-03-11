-- name: GetDocument :one
SELECT * FROM document
WHERE id = $1 LIMIT 1;

-- name: GetDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND validated = true;

-- name: GetDocumentsFromParent :many
SELECT * FROM document
where parent_id = $1 and validated = true;

-- name: GetRootDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND parent_id is NULL;

-- name: GetUnvalidatedDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND validated = false;

-- name: CreateDocument :one
INSERT INTO document (
    parent_id, customer_id, filename, type, size_bytes, sha_256
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: MarkDocumentAsUploaded :one
UPDATE document
SET validated = true
WHERE id = $1
RETURNING *;
