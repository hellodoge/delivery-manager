SELECT ip_address, issued_at, expires_at, invalidated, id
FROM refresh_tokens
WHERE user_id = $1
  AND issued_at > $2;