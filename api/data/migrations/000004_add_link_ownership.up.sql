ALTER TABLE links ADD COLUMN user_id UUID REFERENCES users (id) ON DELETE SET NULL;
ALTER TABLE links ADD COLUMN creator_ip TEXT;

CREATE INDEX idx_links_user_id ON links (user_id);
CREATE INDEX idx_links_creator_ip ON links (creator_ip, created_at);
