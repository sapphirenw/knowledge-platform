-- name: CreateBetaApiKey :one
INSERT INTO beta_api_key ( name, is_admin )
VALUES ( $1, $2 )
RETURNING *;

-- name: GetBetaApiKey :one
SELECT * FROM beta_api_key
WHERE id = $1;