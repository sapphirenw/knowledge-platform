/*
############################################################
Blog Post
############################################################
*/

-- configurations for the generated blog posts
CREATE TABLE blog_post_config(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata
    main_topic TEXT NOT NULL,
    url TEXT,
    metadata JSON NOT NULL DEFAULT '{}',

    -- general configuration
    min_sections INT NOT NULL DEFAULT 4, -- min number of sections to generate
    max_sections INT NOT NULL DEFAULT 10, -- max number of sections to generate
    documents_per_section INT NOT NULL DEFAULT 3, -- number of documents to use as references in sections
    website_pages_per_section INT NOT NULL DEFAULT 3, -- number of pages to use as references in sections

    -- auto generation configuration
    -- auto_gen BOOLEAN NOT NULL DEFAULT FALSE,
    -- auto_gen_cadence TEXT NOT NULL DEFAULT '24h',
    -- auto_gen_time TIME NOT NULL DEFAULT '00:00:00', -- time of the day this content is to be created

    -- llm config
    llm_content_generation_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE CASCADE,
    llm_vector_summarization_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE CASCADE,
    llm_website_summarization_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE CASCADE,
    llm_proof_reading_default_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- categories that get assigned to a blog post
-- single category per blog post
CREATE TABLE blog_category(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    title VARCHAR(20) NOT NULL,
    text_color_hex VARCHAR(7) DEFAULT NULL, -- "#ffffff"
    bg_color_hex VARCHAR(7) DEFAULT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- holds the metadata for a blog post
CREATE TABLE blog_post(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_library_id uuid NOT NULL REFERENCES project_library(id) ON DELETE CASCADE,
    blog_category_id uuid DEFAULT NULL REFERENCES blog_category(id) ON DELETE SET NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL,
    metadata JSON NOT NULL DEFAULT '{}',

    PRIMARY KEY (id),
    CONSTRAINT cnst_blog_post_unique_title UNIQUE (customer_id, title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- section for a blog post outline. Controls how itself is generated.
-- assets are limited to 1 per section
CREATE TABLE blog_post_section(
    id uuid NOT NULL DEFAULT uuid7(),
    blog_post_id uuid NOT NULL REFERENCES blog_post(id) ON DELETE CASCADE,

    additional_instructions TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL,
    asset_id uuid DEFAULT NULL, -- from asset_catalog
    metadata JSON NOT NULL DEFAULT '{}',

    -- models when null, uses customer/defined default
    content_generation_model_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    vector_summarization_model_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    website_summarization_model_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,
    proof_reading_model_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- website pages used as reference for the section
-- can be automatically assigned
CREATE TABLE blog_post_section_website_page(
    id uuid NOT NULL DEFAULT uuid7(),
    blog_post_section_id uuid NOT NULL REFERENCES blog_post_section(id) ON DELETE CASCADE,
    website_page_id uuid NOT NULL REFERENCES website_page(id) ON DELETE CASCADE,
    query TEXT NOT NULL, -- query to use when querying vectorstore

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- documents used as reference for the section
-- can be automatically assigned
CREATE TABLE blog_post_section_document(
    id uuid NOT NULL DEFAULT uuid7(),
    blog_post_section_id uuid NOT NULL REFERENCES blog_post_section(id) ON DELETE CASCADE,
    document_id uuid NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    query TEXT NOT NULL, -- query to use when querying vectorstore

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- section for the content. Stores multiple versions for a section to enable
-- composition of version and feedback as a conversation.
CREATE TABLE blog_post_section_content(
    id uuid NOT NULL DEFAULT uuid7(),
    blog_post_section_id uuid NOT NULL REFERENCES blog_post_section(id) ON DELETE CASCADE,

    content TEXT NOT NULL, -- raw content that the user can edit / give feedback for
    feedback TEXT NOT NULL DEFAULT '', -- feedback is ALWAYS used after the content in the conversation
    index INT NOT NULL, -- index of the conversation

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- tags on a blog post
-- multiple tags per blog post
CREATE TABLE blog_post_tag(
    id uuid NOT NULL DEFAULT uuid7(),
    blog_post_id uuid NOT NULL REFERENCES blog_post(id) ON DELETE CASCADE,

    title VARCHAR(15) NOT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);