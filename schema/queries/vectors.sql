-- name: GetVectorsByCustomer :many
SELECT * FROM vector_store
WHERE customer_id = $1;

-- name: GetVectorsByDocument :many
SELECT * FROM vector_store
WHERE document_id = $1;

-- name: GetVector :one
SELECT * FROM vector_store
WHERE id = $1;

-- name: CreateVector :one
INSERT INTO vector_store (
    raw, embeddings, customer_id, document_id, index
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;