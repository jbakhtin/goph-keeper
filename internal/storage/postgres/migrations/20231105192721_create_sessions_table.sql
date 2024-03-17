-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id bigserial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    refresh_token text,
    finger_print jsonb,
    expire_at timestamp NOT NULL,
    closed_at timestamp,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

ALTER TABLE sessions
    ADD CONSTRAINT user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE TRIGGER set_timestamp_trigger_sessions
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
