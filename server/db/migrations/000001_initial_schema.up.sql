-- DROP DATABASE IF EXISTS cobaltcms;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DATABASE cobaltcms;

CREATE TABLE IF NOT EXISTS public.contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type TEXT NOT NULL,
    name TEXT NOT NULL,
    body JSONB DEFAULT '{}'::jsonb,
    extended_attributes JSONB DEFAULT '{}'::jsonb,
    published_status TEXT NOT NULL,
    publish_start TIMESTAMPTZ,
    publish_end TIMESTAMPTZ,
    created_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.content_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type_name TEXT NOT NULL,
    created_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.published_statuses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status TEXT NOT NULL,
    created_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
