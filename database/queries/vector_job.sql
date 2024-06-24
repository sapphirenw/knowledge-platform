-- name: CreateVectorizeJob :one
INSERT INTO vectorize_job (
    customer_id, documents, websites
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetVectorizeJobsStatus :many
SELECT * FROM vectorize_job
WHERE status = $1;

-- name: GetCustomerVectorizeJobs :many
SELECT * FROM vectorize_job
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: GetVectorizeJob :one
SELECT * FROM vectorize_job
WHERE id = $1;

-- name: UpdateVectorizeJobStatus :one
UPDATE vectorize_job SET
    status = $2,
    message = $3
WHERE id = $1
RETURNING *;

-- name: CreateVectorizeItem :one
INSERT INTO vectorize_item (
    job_id, object_id, error
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetVectorizeJobItems :many
SELECT * FROM vectorize_item
WHERE job_id = $1;