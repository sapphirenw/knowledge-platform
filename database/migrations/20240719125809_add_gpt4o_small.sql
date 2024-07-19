-- +goose Up
-- +goose StatementBegin
-- openai
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'gpt-4o-mini',
    'openai',
    'GPT-4o Mini',
    'A small but powerful model, cheaper and more intelligent compared to gpt-3.5.',
    128000,
    8192,
    0.15,
    0.60
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM available_model WHERE id = 'gpt-4o-mini';
-- +goose StatementEnd
