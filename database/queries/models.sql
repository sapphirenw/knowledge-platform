-- name: GetAvailableModels :many
SELECT * FROM available_model;

-- name: GetAvailableModelsScoped :many
SELECT * FROM available_model
WHERE provider = $1;

-- name: GetAvailableModel :one
SELECT * FROM available_model
WHERE id = $1;