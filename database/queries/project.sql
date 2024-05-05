-- name: CreateProject :one
INSERT INTO project (
    customer_id, title, topic, idea_generation_model_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProjects :many
SELECT * FROM project
WHERE customer_id = $1;

-- name: GetProject :one
SELECT * FROM project
WHERE id = $1;

-- name: CreateProjectIdea :one
INSERT INTO project_idea (
    customer_id, project_id, title
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetProjectIdeas :many
SELECT * FROM project_idea
WHERE project_id = $1;

-- name: GetProjectIdea :one
SELECT * FROM project_idea
WHERE id = $1;

-- name: SetProjectIdeaUsed :one
UPDATE project_idea
    SET used = true
WHERE id = $1
RETURNING *;

