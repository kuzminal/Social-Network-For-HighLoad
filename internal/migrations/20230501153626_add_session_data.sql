-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS social.session
(
    id         serial PRIMARY KEY NOT NULL,
    user_id    varchar UNIQUE     NOT NULL,
    token      varchar            NOT NULL,
    created_at timestamp          NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS social.session
-- +goose StatementEnd
