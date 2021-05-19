DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM products_lists WHERE list_id = $1 AND product_id = $2) THEN
            INSERT INTO products_lists (list_id, product_id, count) VALUES ($1, $2, $3);
        ELSE
            UPDATE products_lists SET count = $3 + count WHERE list_id = $1 AND product_id = $2;
        END IF;
    END
$$