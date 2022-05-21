-- +goose Up
-- +goose StatementBegin
alter table if exists users add column if not exists tg_username text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists users drop column if exists tg_username;
-- +goose StatementEnd
