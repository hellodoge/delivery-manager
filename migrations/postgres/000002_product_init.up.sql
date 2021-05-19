CREATE TABLE products
(
    id          serial       not null unique,
    title       varchar(128) not null,
    description varchar(256),
    price       int          not null default -1
);

CREATE TABLE products_lists
(
    id         serial                                         not null unique,
    product_id int references products (id) on delete cascade not null,
    list_id    int references lists (id) on delete cascade    not null,
    count      int default 1                                  not null
)