DELETE FROM public.content_types
WHERE type_name IN ('DRAFT', 'PUBLISHED', 'SCHEDULED', 'ARCHIVED');
