INSERT INTO users (id, email, password_hash) VALUES
  ('11111111-1111-1111-1111-111111111111', 'seed@example.com', 'seedhash');

INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'seedtokenhash', now() + interval '1 hour');
