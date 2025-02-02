INSERT INTO public.contents
(content_type, name, body, published_status, publish_start, publish_end, created_ts, updated_ts)
VALUES (
  'PAGE', 'Front Page', '{"contentType":"html", "content":"<h1 class\"headline\">Hey Dude!</h1>"}'::jsonb, 'PUBLISHED', null, null, now(), now()
),(
  'PAGE', 'New Page', '{"contentType":"html", "content":"<h1 class\"headline\">New Dude!</h1>"}'::jsonb, 'PUBLISHED', null, null, now(), now()
),(
  'SECTION', 'Sale hero', '{"contentType":"html", "content":"<section class\"sale-section\">Huge sale</section>"}'::jsonb, 'SCHEDULED', '2025-01-01 12:00:00', '2025-01-02 12:00:00', now(), now()
);
