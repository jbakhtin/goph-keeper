-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id bigserial NOT NULL PRIMARY KEY,
   email text NOT NULL UNIQUE,
   password text NOT NULL,
   created_at timestamp NOT NULL DEFAULT now(),
   updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp_trigger_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
