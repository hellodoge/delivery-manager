SELECT *
FROM products
WHERE (NOT $1) AND (
        octet_length($2) > 0 AND position(lower($2) in lower(title)) > 0
        OR octet_length($3) > 0 AND position(lower($3) in lower(description)) > 0
        OR octet_length($4) > 0 AND
           (position(lower($4) in lower(title)) > 0 OR position(lower($4) in lower(description)) > 0)
    )
   OR $1 AND (
        (octet_length($2) = 0 OR position(lower($2) in lower(title)) > 0)
        AND (octet_length($3) = 0 OR position(lower($3) in lower(description)) > 0)
        AND (octet_length($4) = 0 OR
             (position(lower($4) in lower(title)) > 0 OR position(lower($4) in lower(description)) > 0))
    )
