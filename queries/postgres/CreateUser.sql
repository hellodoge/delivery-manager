INSERT INTO users (name, username, password_hash, password_salt)
VALUES ($1, $2, $3, $4)
RETURNING id