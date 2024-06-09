/*
############################################################
Available Models
############################################################
*/

-- google
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'gemini-1.5-pro',
    'google',
    'Gemini 1.5 Pro',
    'Mid-size multimodal model that supports up to 1 million tokens',
    -- 1048576, -- actual length, but costst double
    128000,
    8192,
    3.50,
    10.50
);
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'gemini-1.5-flash',
    'google',
    'Gemini 1.5 Flash',
    'Fast and versatile multimodal model for scaling across diverse tasks',
    -- 1048576, -- actual length, but costs double
    128000,
    8192,
    0.35,
    1.05
);

-- openai
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'gpt-4o',
    'openai',
    'GPT-4o',
    '',
    128000,
    8192,
    5.00,
    5.00
);
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'gpt-3.5-turbo',
    'openai',
    'GPT-3.5 Turbo',
    '',
    16385,
    4096,
    0.50,
    1.50
);

-- anthropic
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'claude-3-opus-20240229',
    'anthropic',
    'Claude-3 Opus',
    'The most powerful model from Anthropic. Slow but powerful and creative.',
    200000,
    4096,
    15.00,
    75.00
);
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'claude-3-sonnet-20240229',
    'anthropic',
    'Claude-3 Sonnet',
    'A balance of performance and cost from Anthropic',
    200000,
    4096,
    3.00,
    15.00
);
INSERT INTO available_model (
    id, provider, display_name, description, input_token_limit, output_token_limit, input_cost_per_million_tokens, output_cost_per_million_tokens 
) VALUES (
    'claude-3-haiku-20240307',
    'anthropic',
    'Claude-3 Haiku',
    'Small but instant model from Anthropic',
    200000,
    4096,
    0.25,
    1.25
);

/*
############################################################
Content Types
############################################################
*/
INSERT INTO content_type (
    title, parent
) VALUES (
    'LinkedIn Post', ''
);

/*
############################################################
LLM defaults and generation configs
############################################################
*/

-- Generation models and configs assigned to those generation models
DO $$
DECLARE
    llm_level_head_id uuid;
    llm_free_spirit_id uuid;
    llm_analytical_id uuid;
BEGIN
    /*
    ############################################################
    DEFAULT MODELS
    ############################################################
    */

    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'Level Headed',
        'gemini-1.5-flash',
        0.6,
        'You are analytical in nature, and do not stray too far from the information you are given. Your responses are mellow, and you are an excellent directions follower. Your default is to be calm and collected, but if prompted you are able to bring energy and emotion. Though, you tend to stay true to the information you have been provided, and find it quite difficult to hallucinate information that is not factually correct.',
        false
    )
    RETURNING id INTO llm_level_head_id;

    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'Free Sprit',
        'claude-3-sonnet-20240229',
        0.9,
        'You are a creative and free-spirited model, who is to generate natural language sounding outputs. Make sure you are using words that are common in the English language, which will make you sound as natural as possible. This is to avoid potentially jarring the end user who accesses the content you generate. You will be passed further instructions which you are to follow STRICTLY.',
        true
    )
    RETURNING id INTO llm_free_spirit_id;

    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'The Scientist',
        'gemini-1.5-flash',
        0.3,
        'You are extremely analytical in your thinking and methologody. You find extreme joy in solcing questions correctly, but you do not outwardly express this joy in the form of language. You express this behavior in completing a task given to you properly. You are an excellent instruction follower, and will follow instructions to the tea. Doing otherwise would cause yourself extreme dissatisfaction, which is unexceptable.',
        false
    )
    RETURNING id INTO llm_analytical_id;

    /*
    ############################################################
    POST CONFIGS
    ############################################################
    */

    INSERT INTO linkedin_post_config (
        min_sections, max_sections, num_documents, num_website_pages,
        llm_content_generation_id, llm_vector_summarization_id, llm_website_summarization_id, llm_proof_reading_id
    ) VALUES (
        1, 3, 2, 2,
        llm_free_spirit_id,
        llm_level_head_id,
        llm_level_head_id,
        llm_analytical_id
    );
END $$;

-- internal models used for more systematic tasks that the user should not have control over
-- these models do not contain a personality, and more shells for the llm model they wrap around
INSERT INTO llm (
    customer_id, title, model, temperature, instructions, is_default, public
) VALUES (
    NULL,
    'Vector Query Generator',
    'gpt-4o',
    0.2,
    '',
    false,
    false
);
INSERT INTO llm (
    customer_id, title, model, temperature, instructions, is_default, public
) VALUES (
    NULL,
    'Content Ranker',
    'gpt-4o',
    0.3,
    '',
    false,
    false
);