INSERT INTO lists (title, description)
VALUES ($1, $2)
RETURNING id