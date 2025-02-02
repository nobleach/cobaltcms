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
SELECT c.id, c.content_type, c.name, c.body, c.extended_attributes
FROM contents_contents cc
LEFT JOIN contents c
ON cc.content_id = c.id
WHERE cc.page_content_id = $1
AND c.published_status = 'PUBLISHED'
OR c.published_status = 'SCHEDULED'
AND  c.publish_end >= TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss')
AND  c.publish_start <=  TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss');
