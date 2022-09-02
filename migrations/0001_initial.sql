-- +migrate Up
CREATE TABLE todos (
    name text
);

-- +migrate Down
DROP TABLE todos;