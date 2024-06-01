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
DO $$
DECLARE
    llm_summarization_id uuid;
    llm_generation_id uuid;
    llm_proof_id uuid;
BEGIN
    /*
    ############################################################
    DEFAULT MODELS
    ############################################################
    */

    /* Sumarization model */
    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'Default Summarization',
        'gpt-3.5-turbo',
        0.5,
        'You are a model that has been specifically designed to summarize content. You are to properly parse the input text, and are to take all relavent facts and defails into account when constructing your summarization. You will be given text as an input, and you will directly reply with the summarization of the content.',
        true
    )
    RETURNING id INTO llm_summarization_id;

    /* Generation model */
    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'Default Generation',
        'claude-3-sonnet-20240229',
        0.9,
        'You are a creative and free-spirited model, who is to generate natural language sounding outputs. Make sure you are using words that are common in the English language, which will make you sound as natural as possible. This is to avoid potentially jarring the end user who accesses the content you generate. You will be passed further instructions which you are to follow STRICTLY.',
        false
    )
    RETURNING id INTO llm_generation_id;

    /* Proof reading model */
    INSERT INTO llm (
        customer_id, title, model, temperature, instructions, is_default
    ) VALUES (
        NULL,
        'Default Proof-reading',
        'claude-3-sonnet-20240229',
        0.3,
        'You are a model that has been crafted to fix mistakes that you see in the outputs/resposnes of humans or other models. Your tasks range from fact checking based on supplied information, spell-checking and document flow, and JSON schema format correction. You are to follow the additional instructions you are given carefully.',
        false
    )
    RETURNING id INTO llm_proof_id;

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
        llm_generation_id,
        llm_summarization_id,
        llm_summarization_id,
        llm_proof_id
    );
END $$;
