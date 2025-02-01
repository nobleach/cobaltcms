-- DROP DATABASE IF EXISTS cobaltcms;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DATABASE cobaltcms;

CREATE TABLE IF NOT EXISTS contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type TEXT NOT NULL,
    name TEXT NOT NULL,
    body TEXT NOT NULL, -- should we do jsonb instead?
    attributes JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
