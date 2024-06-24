-- name: GetDocument :one
SELECT * FROM document
WHERE id = $1 LIMIT 1;

-- name: GetDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND validated = true;

-- name: GetDocumentsFromParent :many
SELECT * FROM document
WHERE parent_id = $1 AND validated = true;

-- name: GetRootDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND parent_id is NULL;

-- name: GetUnvalidatedDocumentsByCustomer :many
SELECT * FROM document
WHERE customer_id = $1 AND validated = false;

-- name: CreateDocument :one
INSERT INTO document (
    parent_id, customer_id, filename, type, size_bytes, sha_256, datastore_type, datastore_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (customer_id, parent_id, filename) DO UPDATE
SET updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: TouchDocument :exec
UPDATE document SET
    updated_at = CURRENT_TIMESTAMP,
    vector_sha_256 = $2
WHERE id = $1;

-- name: UpdateDocumentSummary :one
UPDATE document SET
    summary = $2,
    summary_sha_256 = $3
WHERE id = $1
RETURNING *;

-- name: MarkDocumentAsUploaded :one
UPDATE document
SET validated = true
WHERE id = $1
RETURNING *;

-- name: DeleteDocumentsOlderThan :exec
DELETE FROM document
WHERE customer_id = $1
AND updated_at < $2;

-- name: GetDocumentsOlderThan :many
SELECT * FROM document
WHERE customer_id = $1
AND updated_at < $2;