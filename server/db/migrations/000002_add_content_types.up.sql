INSERT INTO public.content_types
(type_name, created_ts, updated_ts)
VALUES ('DRAFT', now(), now()),
  ('PUBLISHED', now(), now()),
  ('SCHEDULED', now(), now()),
  ('ARCHIVED', now(), now());
