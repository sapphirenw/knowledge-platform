/*
############################################################
LLM default vaules
############################################################
*/
INSERT INTO llm (
    customer_id, title, model, temperature, instructions, is_default
) VALUES (
    NULL,
    'Basic',
    'gpt-3.5-turbo',
    1.00,
    'You are a friendly, AI Assistant here to help and answer all questions politely and concisely.',
    true
);
INSERT INTO llm (
    customer_id, title, model, temperature, instructions, is_default
) VALUES (
    NULL,
    'Direct',
    'gpt-3.5-turbo',
    1.30,
    'You are direct and straight forward. You do not mess around or dilly-dally.',
    true
);