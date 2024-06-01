-- name: CreateLinkedInPost :one
INSERT INTO linkedin_post(
    project_id, project_library_id, project_idea_id, title
) VALUES (
    $1, $2, $3, $4
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

-- name: CreateLinkedInPostConfig :one
INSERT INTO linkedin_post_config (
    project_id, linkedin_post_id,
    min_sections, max_sections, num_documents, num_website_pages,
    llm_content_generation_id, llm_vector_summarization_id, llm_website_summarization_id, llm_proof_reading_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (linkedin_post_id)
DO UPDATE SET
    min_sections = EXCLUDED.min_sections,
    max_sections = EXCLUDED.max_sections,
    num_documents = EXCLUDED.num_documents,
    num_website_pages = EXCLUDED.num_website_pages,
    llm_content_generation_id = EXCLUDED.llm_content_generation_id,
    llm_vector_summarization_id = EXCLUDED.llm_vector_summarization_id,
    llm_website_summarization_id = EXCLUDED.llm_website_summarization_id,
    llm_proof_reading_id = EXCLUDED.llm_proof_reading_id
RETURNING *;

-- name: GetLinkedInPostConfig :one
WITH LinkedInPostConfig AS (
    -- First, try to find a post-specific config
    SELECT * FROM linkedin_post_config
    WHERE linkedin_post_config.project_id = $1 AND linkedin_post_config.linkedin_post_id = $2

    UNION ALL

    -- Fallback to the project's default
    SELECT * FROM linkedin_post_config
    WHERE linkedin_post_config.project_id = $1 AND linkedin_post_config.linkedin_post_id IS NULL

    UNION ALL
    
    -- Finally fallback to global config
    SELECT * FROM linkedin_post_config
    WHERE linkedin_post_config.project_id IS NULL
)
SELECT * FROM LinkedInPostConfig
LIMIT 1;