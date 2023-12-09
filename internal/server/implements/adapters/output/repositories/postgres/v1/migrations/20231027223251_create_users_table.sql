-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id bigserial NOT NULL PRIMARY KEY,
   email text NOT NULL UNIQUE,
   password text NOT NULL,
   created_at timestamp NOT NULL,
   updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
