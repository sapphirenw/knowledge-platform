-- name: CreateLinkedInPost :one
INSERT INTO linkedin_post(
    project_id, project_library_id, project_idea_id, additional_instructions, title
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLinkedInPosts :many
SELECT * FROM linkedin_post
WHERE project_id = $1;

-- name: GetLinkedInPost :one
SELECT * FROM linkedin_post
WHERE id = $1;

-- name: GetLinkedInPostLibrary :one
SELECT * FROM linkedin_post
WHERE project_library_id = $1;