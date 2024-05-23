-- name: GetTokenUsage :many
SELECT * FROM token_usage
WHERE customer_id = $1;

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

