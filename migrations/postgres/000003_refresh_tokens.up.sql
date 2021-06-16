CREATE TABLE refresh_tokens
(
    token varchar(32) primary key not null unique,
    user_id int references users (id) on delete cascade not null,
    ip_address varchar(64),
    issued_at timestamp without time zone not null default now(),
    expires_at timestamp without time zone not null default now() + interval '1 DAY',
    invalidated bool default false
);