INSERT INTO refresh_tokens (token, user_id, ip_address, expires_at)
VALUES ($1, $2, $3, $4);