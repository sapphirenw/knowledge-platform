/*
############################################################
LinkedIn Post
############################################################
*/

-- linkedin posts
CREATE TABLE linkedin_post(
    id uuid NOT NULL DEFAULT uuid7(),
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    project_library_id uuid NOT NULL REFERENCES project_library(id) ON DELETE CASCADE,
    
    project_idea_id uuid NULL REFERENCES project_idea(id) ON DELETE SET NULL,
    
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

-- general config for creating linkedin posts
    CREATE TABLE linkedin_post_config(
        id uuid NOT NULL DEFAULT uuid7(),

        -- null: default for all users
        -- not null: tied to a customer's project
        project_id uuid NULL REFERENCES project(id) ON DELETE CASCADE,

        -- null: default for the entire project
        -- not null: config for the specific post
        linkedin_post_id uuid NULL REFERENCES linkedin_post(id) ON DELETE CASCADE,

        -- general config
        min_sections INT NOT NULL DEFAULT 1,
        max_sections INT NOT NULL DEFAULT 2,
        num_documents INT NOT NULL DEFAULT 2,
        num_website_pages INT NOT NULL DEFAULT 2,

        -- llm config
        llm_content_generation_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
        llm_vector_summarization_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
        llm_website_summarization_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
        llm_proof_reading_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,

        PRIMARY KEY (id),
        CONSTRAINT cnst_unique_linkedin_post_config UNIQUE
        (linkedin_post_id),

        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );