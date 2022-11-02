CREATE TABLE stock
(
    id         uuid    NOT NULL PRIMARY KEY,
    name       varchar NOT NULL,
    quantity   bigint DEFAULT 0,
    created_at timestamp,
    updated_at timestamp
)