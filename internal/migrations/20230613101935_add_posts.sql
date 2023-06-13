-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS social.posts
(
    id             varchar PRIMARY KEY NOT NULL,
    author_user_id varchar             NOT NULL,
    "text"         varchar             NOT NULL,
    created_at     timestamp           NOT NULL
);
CREATE TABLE IF NOT EXISTS social.friends
(
    user_id    varchar   NOT NULL,
    friend_id  varchar   NOT NULL,
    created_at timestamp NOT NULL,
    PRIMARY KEY (user_id, friend_id)
);
CREATE INDEX friends_ind ON social.friends (user_id, friend_id);
CREATE INDEX posts_author_ind ON social.posts (author_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX social.friends_ind;
DROP INDEX social.posts_author_ind;
DROP TABLE IF EXISTS social.posts;
DROP TABLE IF EXISTS social.friends;

-- +goose StatementEnd
