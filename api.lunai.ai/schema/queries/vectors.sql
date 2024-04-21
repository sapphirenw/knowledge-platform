-- name: CreateVector :one
INSERT INTO vector_store (
    raw, embeddings, customer_id
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: CreateDocumentVector :one
INSERT INTO document_vector (
    document_id, vector_store_id, customer_id, index
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: CreateWebsitePageVector :one
INSERT INTO website_page_vector (
    website_page_id, vector_store_id, customer_id, index
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;