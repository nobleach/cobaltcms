# Cobalt CMS

CobaltCMS is built to solve three problems:

- Scheduling of content releases
- Previewing of content that's in DRAFT mode
- Cost of these features

When looking at other solutions, it's evident that they either exist to
make money, or they're base on PHP.

My needs are simple, I want to release blobs of JSON to be consumed by
any front-end application, and do it on a shedule. I also want whomever
is creating that content to see it in a preview mode before it goes live.

## Server

The server is written in Go + Echo. It's a simple REST API that allows
content to be created, read, updated and deleted.

## Client

There is none. Perhaps I'll create one in the future.

## Database

Right now, Postgres is the only database supported. The actual content
is stored in a JSONB column. Metadata about the content is stored in
normal VARCHAR/Integer/DateTime/Timestamp columns.

## Development

## Deployment
