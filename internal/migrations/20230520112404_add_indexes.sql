-- +goose Up
-- +goose StatementBegin
create index search_fname on social.users(first_name text_pattern_ops, second_name text_pattern_ops);
CREATE INDEX token_ind ON social.session (token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX social.search_fname;
DROP INDEX social.token_ind;
-- +goose StatementEnd
