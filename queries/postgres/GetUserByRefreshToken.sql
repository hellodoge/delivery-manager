SELECT users.*
FROM refresh_tokens rt
         INNER JOIN users
                    ON rt.user_id = users.id
WHERE rt.expires_at > now()
  AND rt.invalidated = false
  AND rt.token = $1;