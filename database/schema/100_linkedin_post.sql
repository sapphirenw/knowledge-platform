/*
############################################################
LinkedIn Post
############################################################
*/

-- general config for creating linkedin posts
CREATE TABLE linkedin_post_config(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- general config
    min_sections INT NOT NULL DEFAULT 1,
    max_sections INT NOT NULL DEFAULT 2,
    documents_per_post INT NOT NULL DEFAULT 2,
    website_pages_per_post INT NOT NULL DEFAULT 2,

    -- llm config
    llm_content_generation_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    llm_vector_summarization_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    llm_website_summarization_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    llm_proof_reading_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- linkedin posts
CREATE TABLE linkedin_post(
    id uuid NOT NULL DEFAULT uuid7(),
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    project_library_id uuid NOT NULL REFERENCES project_library(id) ON DELETE CASCADE,
    
    project_idea_id uuid NULL REFERENCES project_idea(id) ON DELETE SET NULL,

    additional_instructions TEXT NOT NULL DEFAULT '',
    
    title TEXT NOT NULL,
    asset_id uuid DEFAULT NULL REFERENCES asset_catalog(id) ON DELETE SET NULL,
    metadata JSON NOT NULL DEFAULT '{}',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- links together the linkedin post with it's conversation used to generate the content
CREATE TABLE linkedin_post_conversation(
    id uuid NOT NULL DEFAULT uuid7(),
    linkedin_post_id uuid NOT NULL REFERENCES linkedin_post(id) ON DELETE CASCADE,
    conversation_id uuid NOT NULL REFERENCES conversation(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);