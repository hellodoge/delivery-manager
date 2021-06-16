SELECT ip_address, issued_at, expires_at, id
FROM refresh_tokens
WHERE user_id = $1
  AND invalidated = false
  AND expires_at > now();