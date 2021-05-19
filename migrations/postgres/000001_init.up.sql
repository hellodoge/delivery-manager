CREATE TABLE users
(
    id serial not null unique,
    name varchar(256) not null,
    username varchar(64) not null unique,
    password_salt varchar(32) not null,
    password_hash varchar(64) not null
);

CREATE TABLE lists
(
    id serial not null unique,
    title varchar(128) not null,
    description varchar(256)
);

CREATE TABLE users_lists
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    list_id int references lists (id) on delete cascade not null
);

CREATE TYPE delivery_status AS ENUM ('Waiting', 'Cancelled', 'Transit', 'Return', 'Done');

CREATE TABLE deliveries
(
    id serial not null unique,
    send_at timestamp with time zone not null,
    arrive_at timestamp with time zone not null,
    status delivery_status not null default 'Waiting',
    list int references lists (id) on delete cascade not null,
    owner int references users (id) on delete cascade not null,
    recipient int references users (id) on delete cascade not null
);
