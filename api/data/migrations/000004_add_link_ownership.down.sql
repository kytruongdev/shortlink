DROP INDEX idx_links_creator_ip;
DROP INDEX idx_links_user_id;

ALTER TABLE links DROP COLUMN creator_ip;
ALTER TABLE links DROP COLUMN user_id;
