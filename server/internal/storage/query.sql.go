// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	types "github.com/nobleach/cobaltcms/internal/types"
)

const getPublishedContentById = `-- name: GetPublishedContentById :many
SELECT c.id, c.fragment_type, c.name, c.body, c.extended_attributes
FROM contents_contents cc
LEFT JOIN contents c
ON cc.content_id = c.id
WHERE cc.page_content_id = $1
AND c.published_status = 'PUBLISHED'
OR c.published_status = 'SCHEDULED'
AND  c.publish_end >= TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss')
AND  c.publish_start <=  TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:ss')
`

type GetPublishedContentByIdParams struct {
	PageContentID uuid.UUID
	ToTimestamp   string
}

type GetPublishedContentByIdRow struct {
	ID                 uuid.NullUUID
	FragmentType       pgtype.Text
	Name               pgtype.Text
	Body               types.JSONB
	ExtendedAttributes types.JSONB
}

func (q *Queries) GetPublishedContentById(ctx context.Context, arg GetPublishedContentByIdParams) ([]GetPublishedContentByIdRow, error) {
	rows, err := q.db.Query(ctx, getPublishedContentById, arg.PageContentID, arg.ToTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPublishedContentByIdRow
	for rows.Next() {
		var i GetPublishedContentByIdRow
		if err := rows.Scan(
			&i.ID,
			&i.FragmentType,
			&i.Name,
			&i.Body,
			&i.ExtendedAttributes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAllPublishedStatuses = `-- name: ListAllPublishedStatuses :many
SELECT id, status
FROM published_statuses
`

type ListAllPublishedStatusesRow struct {
	ID     uuid.UUID
	Status string
}

func (q *Queries) ListAllPublishedStatuses(ctx context.Context) ([]ListAllPublishedStatusesRow, error) {
	rows, err := q.db.Query(ctx, listAllPublishedStatuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAllPublishedStatusesRow
	for rows.Next() {
		var i ListAllPublishedStatusesRow
		if err := rows.Scan(&i.ID, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPublishedContentForDateTime = `-- name: ListPublishedContentForDateTime :many
SELECT id, fragment_type, name, body, extended_attributes 
FROM contents
WHERE published_status = 'PUBLISHED'
OR published_status = 'SCHEDULED'
AND  publish_end >= TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss')
AND  publish_start <=  TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:ss')
`

type ListPublishedContentForDateTimeRow struct {
	ID                 uuid.UUID
	FragmentType       string
	Name               string
	Body               types.JSONB
	ExtendedAttributes types.JSONB
}

func (q *Queries) ListPublishedContentForDateTime(ctx context.Context, toTimestamp string) ([]ListPublishedContentForDateTimeRow, error) {
	rows, err := q.db.Query(ctx, listPublishedContentForDateTime, toTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPublishedContentForDateTimeRow
	for rows.Next() {
		var i ListPublishedContentForDateTimeRow
		if err := rows.Scan(
			&i.ID,
			&i.FragmentType,
			&i.Name,
			&i.Body,
			&i.ExtendedAttributes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveContentRelations = `-- name: SaveContentRelations :execlastid
INSERT INTO contents_contents (page_content_id, content_id, created_ts, updated_ts)
VALUES ($1, $2, now(), now())
`

type SaveContentRelationsParams struct {
	PageContentID uuid.UUID
	ContentID     uuid.UUID
}

const saveNewContent = `-- name: SaveNewContent :one
INSERT INTO contents (fragment_type, name, body, extended_attributes, published_status, publish_start, publish_end, created_ts, updated_ts)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING id
`

type SaveNewContentParams struct {
	FragmentType       string
	Name               string
	Body               types.JSONB
	ExtendedAttributes types.JSONB
	PublishedStatus    string
	PublishStart       *time.Time
	PublishEnd         *time.Time
}

func (q *Queries) SaveNewContent(ctx context.Context, arg SaveNewContentParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, saveNewContent,
		arg.FragmentType,
		arg.Name,
		arg.Body,
		arg.ExtendedAttributes,
		arg.PublishedStatus,
		arg.PublishStart,
		arg.PublishEnd,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
