-- name: CreateTokenUsage :one
INSERT INTO token_usage (
    id, customer_id, conversation_id, model, input_tokens, output_tokens, total_tokens
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (id)
DO UPDATE SET
    customer_id = EXCLUDED.customer_id,
    conversation_id = EXCLUDED.conversation_id,
    model = EXCLUDED.model,
    input_tokens = EXCLUDED.input_tokens,
    output_tokens = EXCLUDED.output_tokens,
    total_tokens = EXCLUDED.total_tokens
RETURNING *;

-- name: GetCustomerTokenUsage :many
SELECT * FROM token_usage
WHERE customer_id = $1
  AND (model = $2 OR $2 = '')
ORDER BY created_at DESC, id DESC
LIMIT $3 OFFSET ($4::int - 1) * $3;

-- name: GetCustomerTokensPageCount :one
SELECT
    CEIL(COUNT(*)::FLOAT / $3) as max_pages
FROM token_usage
WHERE customer_id = $1
AND (model = $2 OR $2 = '');

-- name: GetCustomerUsageGrouped :many
SELECT
    tu.model AS model,
    SUM(tu.input_tokens) AS input_tokens_sum,
    SUM(tu.output_tokens) AS output_tokens_sum,
    SUM(tu.total_tokens) AS total_tokens_sum,
    am.input_cost_per_million_tokens,
    am.output_cost_per_million_tokens
FROM token_usage tu
JOIN available_model am ON tu.model = am.id
WHERE tu.customer_id = $1
GROUP BY tu.model, am.input_cost_per_million_tokens, am.output_cost_per_million_tokens;
