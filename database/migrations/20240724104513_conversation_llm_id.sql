-- +goose Up
-- +goose StatementBegin
ALTER TABLE conversation ADD COLUMN curr_llm_id uuid NULL REFERENCES llm(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE conversation DROP COLUMN curr_llm_id;
-- +goose StatementEnd
