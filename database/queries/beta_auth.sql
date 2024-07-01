-- name: CreateBetaApiKey :one
INSERT INTO beta_api_key ( name )
VALUES ( $1 )
RETURNING *;

-- name: GetBetaApiKey :one
SELECT * FROM beta_api_key
WHERE id = $1;