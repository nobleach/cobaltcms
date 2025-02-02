-- name: ListAllPublishedStatuses :many
SELECT id, status
FROM published_statuses;

-- name: ListPublishedContentForDateTime :many
SELECT id, content_type, name, body, extended_attributes 
FROM contents
WHERE published_status = 'PUBLISHED'
OR published_status = 'SCHEDULED'
AND  publish_end >= TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss')
AND  publish_start <=  TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss');

-- name: GetPublishedContentById :many
SELECT id, content_type, name, body, extended_attributes 
FROM contents
WHERE id = $1
OR parent_id = $1
AND published_status = 'PUBLISHED'
OR published_status = 'SCHEDULED'
AND  publish_end >= TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss')
AND  publish_start <=  TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss');
