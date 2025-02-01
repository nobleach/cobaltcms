INSERT INTO public.published_statuses
(status, created_ts, updated_ts)
VALUES ('DRAFT', now(), now()),
  ('PUBLISHED', now(), now()),
  ('SCHEDULED', now(), now()),
  ('ARCHIVED', now(), now());
