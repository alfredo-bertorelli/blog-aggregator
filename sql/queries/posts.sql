-- name: CreatePost :one
INSERT INTO posts (id, title, url, description, published_at, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetRecentPostsForUser :many
SELECT posts.*
FROM posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id=$1
ORDER BY COALESCE(posts.published_at, posts.created_at) DESC
LIMIT $2;
