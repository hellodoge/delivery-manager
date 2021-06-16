UPDATE refresh_tokens
SET invalidated = true
WHERE id = $1
  AND user_id = $2;