-- name: GetWebsite :one
SELECT * FROM website
WHERE id = $1;

-- name: GetWebsitesByCustomer :many
SELECT * FROM website
WHERE customer_id = $1;

-- name: CreateWebsite :one
INSERT INTO website (
    customer_id, protocol, domain, blacklist, whitelist
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: CreateWebsitePage :one
INSERT INTO website_page (
    customer_id, website_id, url, sha_256
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateWebsitePageSignature :one
UPDATE website_page SET
    sha_256 = $2
WHERE id = $1
RETURNING *;

-- name: GetWebsitePagesBySite :many
SELECT * FROM website_page
WHERE website_id = $1;