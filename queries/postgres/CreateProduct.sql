INSERT INTO products (title, description, price)
VALUES ($1, $2, $3)
RETURNING id