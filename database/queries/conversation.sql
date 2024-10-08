-- name: CreateConversation :one
INSERT INTO conversation (
    customer_id, title, conversation_type, system_message, metadata
) VALUES ( $1, $2, $3, $4, $5 )
RETURNING *;

-- name: GetConversations :many
SELECT * FROM conversation
WHERE customer_id = $1
ORDER BY updated_at DESC;

-- name: GetConversationsWithCount :many
SELECT c.*, COUNT(cm.id) AS message_count
FROM conversation c
JOIN conversation_message cm
ON c.id = cm.conversation_id
WHERE c.customer_id = $1
GROUP BY c.id
ORDER BY c.updated_at DESC;

-- name: GetConversation :one
SELECT * FROM conversation
WHERE id = $1;

-- name: SetConversationError :one
UPDATE conversation SET
    has_error = true,
    error_message = $2
WHERE id = $1
RETURNING *;

-- name: CreateConversationMessage :one
INSERT INTO conversation_message (
    conversation_id,
    llm_id,
    model,
    temperature,
    instructions,
    role,
    message,
    index,
    tool_use_id,
    tool_name,
    tool_arguments,
    tool_results
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 )
ON CONFLICT (conversation_id, index)
DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetConversationMessages :many
SELECT * FROM conversation_message
WHERE conversation_id = $1
ORDER BY index ASC;

-- name: ClearConversation :exec
DELETE FROM conversation_message
WHERE conversation_id = $1;

-- name: UpdateConversationTitle :one
UPDATE conversation SET
    title = $2
WHERE id = $1
RETURNING *;

-- name: SetChatLLM :exec
UPDATE conversation SET
    curr_llm_id = $2
WHERE id = $1;

-- name: GetChatLLM :one
SELECT 
    sqlc.embed(llm),
    sqlc.embed(am)
FROM conversation c
JOIN llm ON llm.id = c.curr_llm_id
JOIN available_model am ON llm.model = am.id
WHERE c.customer_id = $1
AND c.id = $2;