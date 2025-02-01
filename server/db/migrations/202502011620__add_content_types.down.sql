DELETE FROM public.content_types
WHERE type_name IN
('PAGE', 'SECTION', 'COMPONENT');
