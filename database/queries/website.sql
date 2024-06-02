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
ON CONFLICT ON CONSTRAINT cnst_unique_website
DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,
    blacklist = EXCLUDED.blacklist,
    whitelist = EXCLUDED.whitelist
RETURNING *;

-- name: CreateWebsitePage :one
INSERT INTO website_page (
    customer_id, website_id, url, sha_256, metadata
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT ON CONSTRAINT cnst_unique_website_page
DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,
    is_valid = TRUE
RETURNING *;

-- name: UpdateWebsitePageSignature :one
UPDATE website_page SET
    sha_256 = $2
WHERE id = $1
RETURNING *;

-- name: GetWebsitePagesBySite :many
SELECT * FROM website_page
WHERE website_id = $1;

-- name: DeleteWebsitePagesOlderThan :exec
DELETE FROM website_page
WHERE customer_id = $1
AND updated_at < $2;

-- name: DeleteWebsitePagesNotValid :exec
DELETE FROM website_page
WHERE customer_id = $1
AND website_id = $2
AND is_valid = FALSE;

-- name: SetWebsitePagesNotValid :exec
UPDATE website_page SET is_valid = FALSE
WHERE customer_id = $1
AND website_id = $2;