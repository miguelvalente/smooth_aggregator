-- name: CreatePost :one
insert into posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostsByUserId :many
select p.*
from posts p
join feeds_follows ff on p.feed_id = ff.feed_id
where ff.user_id = $1;
