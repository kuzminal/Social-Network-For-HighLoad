-- +goose Up
-- +goose StatementBegin
create table IF NOT EXISTS social.messages
(
    id         varchar,
    from_user  varchar   not null,
    text       varchar   not null,
    to_user    varchar   not null,
    chat_id    varchar   not null,
    created_at timestamp not null,
    constraint messages_pk
        primary key (id, chat_id)
);

CREATE TABLE IF NOT EXISTS social.chats
(
    chat_id   varchar,
    user_from varchar,
    user_to   varchar,
    CONSTRAINT chats_pkey PRIMARY KEY (chat_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS social.messages;
DROP TABLE IF EXISTS social.chats;
-- +goose StatementEnd
