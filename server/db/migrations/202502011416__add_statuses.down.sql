DELETE FROM public.published_statuses
WHERE status IN ('DRAFT', 'PUBLISHED', 'SCHEDULED', 'ARCHIVED');
