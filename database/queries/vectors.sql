-- name: CreateVector :one
INSERT INTO vector_store (
    customer_id, raw, embeddings, content_type, object_id, object_parent_id, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: CreateDocumentVector :one
INSERT INTO document_vector (
    document_id, vector_store_id, customer_id, index, metadata
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListDocumentVectors :many
SELECT * FROM document_vector
WHERE customer_id = $1;

-- name: CreateWebsitePageVector :one
INSERT INTO website_page_vector (
    website_page_id, vector_store_id, customer_id, index, metadata
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListWebsitePageVectors :many
SELECT * FROM website_page_vector
WHERE customer_id = $1;

-- name: QueryVectorStoreRaw :many
SELECT * FROM vector_store
WHERE customer_id = $1
ORDER BY embeddings <#> $3
LIMIT $2;

-- name: QueryVectorStoreDocuments :many
SELECT d.*
FROM vector_store vs
JOIN document_vector dv ON vs.id = dv.vector_store_id
JOIN document d ON d.id = dv.document_id
WHERE vs.customer_id = $1
ORDER BY vs.embeddings <#> $3
LIMIT $2;

-- name: QueryVectorStoreDocumentsScoped :many
SELECT d.*
FROM vector_store vs
JOIN document_vector dv ON vs.id = dv.vector_store_id
JOIN document d ON d.id = dv.document_id
JOIN folder f ON f.id = d.parent_id
WHERE vs.customer_id = $1
AND (f.id = ANY($4::uuid[]) OR d.id = ANY($5::uuid[]))
ORDER BY vs.embeddings <#> $3
LIMIT $2;

-- name: QueryVectorStoreWebsitePages :many
SELECT wp.*
FROM vector_store vs
JOIN website_page_vector wpv ON vs.id = wpv.vector_store_id
JOIN website_page wp ON wp.id = wpv.website_page_id
WHERE vs.customer_id = $1
ORDER BY vs.embeddings <#> $3
LIMIT $2;

-- name: QueryVectorStoreWebsitePagesScoped :many
SELECT wp.*
FROM vector_store vs
JOIN website_page_vector wpv ON vs.id = wpv.vector_store_id
JOIN website_page wp ON wp.id = wpv.website_page_id
JOIN website w ON w.id = wp.website_id
WHERE vs.customer_id = $1
AND (w.id = ANY($4::uuid[]) OR wp.id = ANY($5::uuid[]))
ORDER BY vs.embeddings <#> $3
LIMIT $2;

-- name: QueryVectorStore :many
SELECT
    vs.*, d.*, wp.*
FROM vector_store vs
LEFT JOIN document_vector dv ON dv.vector_store_id = vs.object_id
LEFT JOIN document d ON d.id = dv.document_id
LEFT JOIN folder f ON f.id = d.parent_id 
LEFT JOIN website_page_vector wpv ON wpv.vector_store_id = vs.object_id
LEFT JOIN website_page wp ON wp.id = wpv.website_page_id
LEFT JOIN website w ON wp.website_id = w.id
WHERE vs.customer_id = $1
AND (
    (d.id = ANY($4::uuid[]) OR $4 IS NULL)
    OR
    (f.id = ANY($5::uuid[]) OR $5 IS NULL)
    OR
    (wp.id = ANY($6::uuid[]) OR $6 IS NULL)
    OR
    (w.id = ANY($7::uuid[]) OR $7 IS NULL)
)
ORDER BY vs.embeddings <#> $3
LIMIT $2;