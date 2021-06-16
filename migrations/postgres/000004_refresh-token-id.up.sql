ALTER TABLE refresh_tokens
ADD COLUMN id serial;

CREATE INDEX refresh_tokens_id_index ON refresh_tokens(id);