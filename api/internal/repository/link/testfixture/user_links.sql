INSERT INTO users (id, email, password_hash) VALUES
  ('22222222-2222-2222-2222-222222222222', 'owner@example.com', 'seedhash');

INSERT INTO links (code, original_url, normalized_url, user_id, created_at) VALUES
  ('older1', 'https://a.com', 'https://a.com', '22222222-2222-2222-2222-222222222222', '2024-01-01T00:00:00Z'),
  ('newer1', 'https://b.com', 'https://b.com', '22222222-2222-2222-2222-222222222222', '2024-01-02T00:00:00Z');
