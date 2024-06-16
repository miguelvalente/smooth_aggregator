-- name: CreateFeedFollows :one
INSERT INTO feeds_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteUserFeedFollows :exec
DELETE FROM feeds_follows f
WHERE f.id = $1;

-- name: GetFeedsFollowsByUserId :many
select * from feeds_follows f
where f.user_id = $1;
