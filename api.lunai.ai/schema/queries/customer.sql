-- name: GetCustomer :one
SELECT * FROM customer
WHERE id = $1 LIMIT 1;

-- name: GetCustomerByName :one
SELECT * FROM customer
WHERE name = $1 LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customer
ORDER BY name;

-- name: CreateCustomer :one
INSERT INTO customer (
    name
) VALUES (
    $1
)
RETURNING *;

-- name: CreateCustomerTest :one
INSERT INTO customer (
    id, name
) VALUES (
    $1, $2
)
ON CONFLICT (id) DO UPDATE
SET updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: UpdateCustomer :exec
UPDATE customer
    set name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customer
WHERE id = $1;
