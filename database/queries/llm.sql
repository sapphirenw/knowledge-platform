-- name: CreateLLM :one
INSERT INTO llm (
    customer_id, title, model, temperature, instructions, is_default
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetLLM :one
SELECT * FROM llm
WHERE id = $1;

-- name: GetDefaultLLM :one
WITH RequiredLLM AS (
    -- First, try to find a customer-specific default if customer_id is provided
    SELECT * FROM llm
    WHERE llm.customer_id = $1 AND llm.is_default = true

    UNION ALL

    -- Fallback to a global default if no customer-specific default is found
    SELECT * FROM llm
    WHERE llm.customer_id IS NULL AND llm.is_default = true
    AND NOT EXISTS (
        SELECT 1
        FROM llm
        WHERE llm.customer_id = $1 AND llm.is_default = true
    )
    
)
SELECT * FROM RequiredLLM
LIMIT 1;

-- name: GetLLMsByCustomer :many
SELECT * FROM llm
WHERE customer_id = $1;

-- name: GetStandardLLMs :many
SELECT * FROM llm
where customer_id IS NULL;