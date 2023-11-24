-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id bigserial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    refresh_token text,
    finger_print jsonb,
    expire_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    closed_at timestamp,
    updated_at timestamp
);

ALTER TABLE sessions
    ADD CONSTRAINT user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
