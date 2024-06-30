-- name: CreateBetaApiKey :one
INSERT INTO beta_api_key DEFAULT VALUES
RETURNING *;

-- name: GetBetaApiKey :one
SELECT * FROM beta_api_key
WHERE id = $1;