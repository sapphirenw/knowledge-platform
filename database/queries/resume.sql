-- name: CreateResume :one
INSERT INTO resume (
    customer_id, title
) VALUES ( $1, $2 )
RETURNING *;

-- name: CreateCustomerResume :one
INSERT INTO resume (
    id, customer_id, title
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: SetResumeTitle :one
UPDATE resume SET
    title = $2
WHERE id = $1
RETURNING *;

-- name: GetResumesCustomer :many
SELECT * FROM resume
WHERE customer_id = $1;

-- name: GetResume :one
SELECT * FROM resume
WHERE id = $1;

-- name: CreateResumeDocument :one
INSERT INTO resume_document (
    resume_id, document_id, is_resume
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetResumeDocuments :many
SELECT d.* FROM resume_document rd
JOIN document d ON d.id = rd.document_id
WHERE rd.resume_id = $1;

-- name: CreateResumeWebsite :one
INSERT INTO resume_website (
    resume_id, website_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetResumeWebsites :many
SELECT w.* FROM resume_website rw
JOIN website w ON w.id = rw.website_id
WHERE rw.resume_id = $1;

-- name: CreateResumeWebsitePage :one
INSERT INTO resume_website_page (
    resume_id, website_page_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetResumeWebsitePages :many
SELECT wp.* FROM resume_website_page rwp
JOIN website_page wp ON wp.id = rwp.website_page_id
WHERE rwp.resume_id = $1;

-- name: CreateResumeAbout :one
INSERT INTO resume_about (
    resume_id, name, email, phone,
    title, location, github, linkedin
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (resume_id)
DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    title = EXCLUDED.title,
    location = EXCLUDED.location,
    github = EXCLUDED.github,
    linkedin = EXCLUDED.linkedin,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetResumeAbout :one
SELECT * FROM resume_about
WHERE resume_id = $1;

-- name: CreateResumeWorkExperience :one
INSERT INTO resume_work_experience (
    resume_id, company, position, location,
    start_date, end_date, is_current, index,
    information
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetResumeWorkExperiences :many
SELECT * FROM resume_work_experience
WHERE resume_id = $1;

-- name: GetResumeWorkExperience :one
SELECT * FROM resume_work_experience
WHERE id = $1;

-- name: UpdateResumeWorkExperience :one
UPDATE resume_work_experience SET
    company = $2,
    position = $3,
    location = $4,
    start_date = $5,
    end_date = $6,
    is_current = $7,
    index = $8
WHERE id = $1
RETURNING *;

-- name: SetResumeWorkExperienceInfo :one
UPDATE resume_work_experience SET
    information = $2
WHERE id = $1
RETURNING *;

-- name: DeleteResumeWorkExperience :exec
DELETE FROM resume_work_experience
WHERE id = $1;

-- name: CreateResumeProject :one
INSERT INTO resume_project (
    resume_id, title, subtitle, link,
    start_date, end_date, information
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: CreateResumeEducation :one
INSERT INTO resume_education (
    resume_id, institution, major, level, gpa,
    location, start_date, end_date, is_current,
    information
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: CreateResumeSkill :one
INSERT INTO resume_skill (
    resume_id, title, items
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetResumeApplications :many
SELECT * FROM resume_application
WHERE resume_id = $1
ORDER BY status;

-- name: CreateResumeApplication :one
INSERT INTO resume_application (
    resume_id, title, link, company_site, raw_text
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;