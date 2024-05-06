-- name: CreateConversation :one
INSERT INTO conversation (
    customer_id, title
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetConversations :many
SELECT * FROM conversation
WHERE customer_id = $1;

-- name: GetConversation :one
SELECT * FROM conversation
WHERE id = $1;

-- name: CreateConversationMessage :one
INSERT INTO conversation_message (
    conversation_id,
    llm_id,
    model,
    temperature,
    instructions,
    role,
    message,
    index
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
ON CONFLICT (conversation_id, index)
DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetConversationMessages :many
SELECT * FROM conversation_message
WHERE conversation_id = $1
ORDER BY index ASC;