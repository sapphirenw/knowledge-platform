-- name: GetWebsite :one
SELECT * FROM website
WHERE id = $1;

-- name: GetWebsitesByCustomer :many
SELECT * FROM website
WHERE customer_id = $1;

-- name: GetWebsitesByCustomerWithCount :many
SELECT w.*, count(wp.*) as page_count FROM website w
JOIN website_page wp ON w.id = wp.website_id
WHERE w.customer_id = $1
GROUP BY w.id;

-- name: CreateWebsite :one
INSERT INTO website (
    customer_id, protocol, domain, path, blacklist, whitelist
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT ON CONSTRAINT cnst_unique_website
DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,
    blacklist = EXCLUDED.blacklist,
    whitelist = EXCLUDED.whitelist
RETURNING *;

-- name: DeleteWebsiteEmpty :exec
DELETE FROM website w
WHERE w.customer_id = $1
AND NOT EXISTS (
    SELECT 1
    FROM website_page wp
    WHERE wp.website_id = w.id
);

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
    updated_at = CURRENT_TIMESTAMP,
    sha_256 = $2
WHERE id = $1
RETURNING *;

-- name: UpdateWebsitePageSummary :one
UPDATE website_page SET
    updated_at = CURRENT_TIMESTAMP,
    summary = $2,
    summary_sha_256 = $3
WHERE id = $1
RETURNING *;

-- name: UpdateWebsitePageVectorSig :exec
UPDATE website_page SET
    updated_at = CURRENT_TIMESTAMP,
    vector_sha_256 = $2
WHERE id = $1;

-- name: TouchWebsitePage :exec
UPDATE website_page SET
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

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
UPDATE website_page SET
    updated_at = CURRENT_TIMESTAMP,
    is_valid = FALSE
WHERE customer_id = $1
AND website_id = $2;