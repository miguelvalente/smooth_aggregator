-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
select * from feeds f;

-- name: GetFeedsbyID :one
select * from feeds f
where f.id = $1;
