-- +goose Up
-- +goose StatementBegin
ALTER TABLE available_model ADD COLUMN is_visible BOOLEAN NOT NULL DEFAULT true;
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens, is_visible
) VALUES (
    'text-embedding-3-small',
    'openai',
    'Text Embeddings 3 Small',
    '',
    8191,
    8191,
    0.02,
    0,
    false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM available_model WHERE id = 'text-embedding-3-small';
ALTER TABLE available_model DROP COLUMN is_visible;
-- +goose StatementEnd
