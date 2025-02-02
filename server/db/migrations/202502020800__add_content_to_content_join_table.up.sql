-- Many to Many table to join sub-contents to a page id
CREATE TABLE IF NOT EXISTS public.contents_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_content_id UUID NOT NULL,
    content_id UUID NOT NULL,
    created_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
