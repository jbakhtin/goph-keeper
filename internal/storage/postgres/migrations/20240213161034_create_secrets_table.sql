-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets (
     id bigserial NOT NULL PRIMARY KEY,
     user_id bigint NOT NULL,
     description text NOT NULL,
     data jsonb NOT NULL,
     created_at timestamp NOT NULL DEFAULT now(),
     updated_at timestamp NOT NULL DEFAULT now()
);

ALTER TABLE secrets
    ADD CONSTRAINT user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE TRIGGER set_timestamp_trigger_secrets
BEFORE UPDATE ON secrets
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
