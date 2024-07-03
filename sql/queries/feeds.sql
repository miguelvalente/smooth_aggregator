-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, last_fetched_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
select * from feeds f;

-- name: GetFeedsbyID :one
select * from feeds f
where f.id = $1;

-- name: GetNextNFeedsToFetch :many
select * from feeds f
order by
    case
        when last_fetched_at is null then 0
        else 1
    end,
    last_fetched_at
    limit $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $2, updated_at=$2
where feeds.id = $1;
