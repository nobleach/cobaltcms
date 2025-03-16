-- name: ListAllPublishedStatuses :many
SELECT id, status
FROM published_statuses;

-- name: ListPublishedContentForDateTime :many
SELECT id, fragment_type, name, body, extended_attributes 
FROM contents
WHERE published_status = 'PUBLISHED'
OR published_status = 'SCHEDULED'
AND  publish_end >= TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss')
AND  publish_start <=  TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss');

-- name: GetPublishedContentById :many
SELECT c.id, c.fragment_type, c.name, c.body, c.extended_attributes
FROM contents_contents cc
LEFT JOIN contents c
ON cc.content_id = c.id
WHERE cc.page_content_id = $1
AND c.published_status = 'PUBLISHED'
OR c.published_status = 'SCHEDULED'
AND  c.publish_end >= TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss')
AND  c.publish_start <=  TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss');

-- name: SaveNewContent :one
INSERT INTO contents (fragment_type, name, body, extended_attributes, published_status, publish_start, publish_end, created_ts, updated_ts)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING id;

-- name: SaveContentRelations :execlastid
INSERT INTO contents_contents (page_content_id, content_id, created_ts, updated_ts)
VALUES ($1, $2, now(), now());
