/*
############################################################
Blog Post
############################################################
*/

-- root blog that a customer can create
-- customers can create mutliple blogs
-- contains generation configuration inside as well
CREATE TABLE blog(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,

    -- metadata
    title TEXT NOT NULL,
    main_topic TEXT NOT NULL,
    url TEXT,
    metadata JSON NOT NULL DEFAULT '{}',

    -- general configuration
    min_sections INT NOT NULL DEFAULT 4, -- min number of sections to generate
    max_sections INT NOT NULL DEFAULT 10, -- max number of sections to generate
    documents_per_section INT NOT NULL DEFAULT 3, -- number of documents to use as references in sections
    website_pages_per_section INT NOT NULL DEFAULT 3, -- number of pages to use as references in sections

    -- auto generation configuration
    auto_gen BOOLEAN NOT NULL DEFAULT FALSE,
    auto_gen_cadence TEXT NOT NULL DEFAULT '24h',
    auto_gen_time TIME NOT NULL DEFAULT '00:00:00', -- time of the day this content is to be created

    -- llm config
    llm_content_generation_default_id uuid DEFAULT NULL,
    llm_vector_summarization_default_id uuid DEFAULT NULL,
    llm_website_summarization_default_id uuid DEFAULT NULL,
    llm_proof_reading_default_id uuid DEFAULT NULL,

    -- constrains
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_generation_model_id FOREIGN KEY

    (llm_content_generation_default_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_vector_summarization_model_id FOREIGN KEY
    (llm_vector_summarization_default_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_website_summarization_model_id FOREIGN KEY
    (llm_website_summarization_default_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_proof_reading_model_id FOREIGN KEY
    (llm_proof_reading_default_id) REFERENCES llm(id) ON DELETE SET NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- entire websites assigned to a blog for reference
-- if none, then all websites in the store can be used
CREATE TABLE blog_reference_website(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    website_id uuid NOT NULL,

    -- constraints
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_reference_website_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_reference_website_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_reference_website_website_id FOREIGN KEY
    (website_id) REFERENCES website(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- entire folders used as reference for the section
-- if none, then all documents in the store can be used
CREATE TABLE blog_reference_folder(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    folder_id uuid NOT NULL,

    -- constraints
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_reference_folder_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_reference_folder_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_reference_folder_folder_id FOREIGN KEY
    (folder_id) REFERENCES folder(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- categories that get assigned to a blog post
-- single category per blog post
CREATE TABLE blog_category(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,

    title VARCHAR(20) NOT NULL,
    text_color_hex VARCHAR(7) DEFAULT NULL, -- "#ffffff"
    bg_color_hex VARCHAR(7) DEFAULT NULL,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_category_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_category_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- to hold query ideas for generating posts in this blog
CREATE TABLE blog_post_idea(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,

    title TEXT NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_idea_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_idea_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- holds the metadata for a blog post
CREATE TABLE blog_post(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_idea_id uuid DEFAULT NULL,
    blog_category_id uuid DEFAULT NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL,
    metadata JSON NOT NULL DEFAULT '{}',

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_blog_post_idea_id FOREIGN KEY
    (blog_post_idea_id) REFERENCES blog_post_idea(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_post_blog_blog_category_id FOREIGN KEY
    (blog_category_id) REFERENCES blog_category(id) ON DELETE SET NULL,

    -- uniques
    CONSTRAINT cnst_blog_post_unique_title UNIQUE
    (customer_id, title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- section for a blog post outline. Controls how itself is generated.
-- assets are limited to 1 per section
CREATE TABLE blog_post_section(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_id uuid NOT NULL,

    -- data
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    additional_instructions TEXT NOT NULL,
    asset_id uuid DEFAULT NULL, -- from asset_catalog
    metadata JSON NOT NULL DEFAULT '{}',

    -- models when null, uses customer/defined default
    content_generation_model_id uuid DEFAULT NULL,
    vector_summarization_model_id uuid DEFAULT NULL,
    website_summarization_model_id uuid DEFAULT NULL,
    proof_reading_model_id uuid DEFAULT NULL,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_section_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_asset_id FOREIGN KEY
    (asset_id) REFERENCES asset_catalog(id) ON DELETE SET NULL,

    CONSTRAINT fk_blog_post_section_generation_model_id FOREIGN KEY
    (content_generation_model_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_post_section_vector_summarization_model_id FOREIGN KEY
    (vector_summarization_model_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_post_section_website_summarization_model_id FOREIGN KEY
    (website_summarization_model_id) REFERENCES llm(id) ON DELETE SET NULL,
    CONSTRAINT fk_blog_post_section_proof_reading_model_id FOREIGN KEY
    (proof_reading_model_id) REFERENCES llm(id) ON DELETE SET NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- website pages used as reference for the section
-- can be automatically assigned
CREATE TABLE blog_post_section_website_page(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_id uuid NOT NULL,
    blog_post_section_id uuid NOT NULL,
    website_page_id uuid NOT NULL,
    query TEXT NOT NULL, -- query to use when querying vectorstore

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_section_website_page_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_website_page_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_website_page_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_website_page_blog_post_section_id FOREIGN KEY
    (blog_post_section_id) REFERENCES blog_post_section(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_website_page_website_page_id FOREIGN KEY
    (website_page_id) REFERENCES website_page(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- documents used as reference for the section
-- can be automatically assigned
CREATE TABLE blog_post_section_document(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_id uuid NOT NULL,
    document_id uuid NOT NULL,
    query TEXT NOT NULL, -- query to use when querying vectorstore

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_section_document_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_document_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_document_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_document_document_id FOREIGN KEY
    (document_id) REFERENCES document(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- section for the content. Stores multiple versions for a section to enable
-- composition of version and feedback as a conversation.
CREATE TABLE blog_post_section_content(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_id uuid NOT NULL,
    blog_post_section_id uuid NOT NULL,

    content TEXT NOT NULL, -- raw content that the user can edit / give feedback for
    feedback TEXT NOT NULL DEFAULT '', -- feedback is ALWAYS used after the content in the conversation
    index INT NOT NULL, -- index of the conversation

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_post_section_content_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_content_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_content_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_post_section_content_blog_post_section_id FOREIGN KEY
    (blog_post_section_id) REFERENCES blog_post_section(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- tags on a blog post
-- multiple tags per blog post
CREATE TABLE blog_post_tag(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    blog_id uuid NOT NULL,
    blog_post_id uuid NOT NULL,

    title VARCHAR(15) NOT NULL,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_tag_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_tag_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_tag_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);