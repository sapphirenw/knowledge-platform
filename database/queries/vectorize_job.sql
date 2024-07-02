-- name: CreateVectorizeJob :one
INSERT INTO vectorize_job (
    customer_id, documents, websites
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetVectorizeJobsWaiting :many
SELECT * FROM vectorize_job vj
WHERE NOT EXISTS (
    SELECT 1
    FROM vectorize_job_item vji
    WHERE vj.id = vji.job_id
);

-- name: GetVectorizeJob :one
SELECT vj.*, vji.status, vji.message, vji.error
FROM vectorize_job vj
JOIN vectorize_job_item vji ON vj.id = vji.job_id
WHERE vj.id = $1
ORDER BY vji.created_at DESC
LIMIT 1;

-- name: GetCustomerVectorizeJobs :many
WITH latest_vji AS (
    SELECT 
        vji.*, 
        ROW_NUMBER() OVER (PARTITION BY vji.job_id ORDER BY vji.created_at DESC) AS rn
    FROM 
        vectorize_job_item vji
)
SELECT 
    vj.*, 
    vji.status, 
    vji.message, 
    vji.error
FROM 
    vectorize_job vj
LEFT JOIN 
    latest_vji vji 
    ON vj.id = vji.job_id 
    AND vji.rn = 1
WHERE 
    vj.customer_id = $1
ORDER BY 
    vj.created_at DESC;

-- name: CreateVectorizeJobItem :one
INSERT INTO vectorize_job_item (
    job_id, status, message, error
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetVectorizeJobItems :many
SELECT * FROM vectorize_job_item
WHERE job_id = $1;