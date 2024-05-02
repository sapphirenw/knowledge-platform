
CREATE EXTENSION IF NOT EXISTS vector;

/*
############################################################
Base Schema
############################################################
*/

CREATE TABLE customer(
    id BIGSERIAL,
    name TEXT NOT NULL,
    datastore VARCHAR(256) NOT NULL DEFAULT 's3', -- name of the datastore the user wants to store their documents

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- generic table to hold vectors for all sorts of data
CREATE TABLE vector_store(
    id BIGSERIAL,
    raw TEXT NOT NULL, -- string utf-8 representation of the data 
    embeddings VECTOR(512) NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    PRIMARY KEY (id, customer_id), -- customer_id needs to exist in the key for partitioning

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) PARTITION BY LIST(customer_id);
CREATE INDEX ON vector_store USING hnsw (embeddings vector_ip_ops);
CREATE TABLE vector_store_default PARTITION OF vector_store DEFAULT; -- default

-- track token usage for a customer across multiple different models
CREATE TABLE token_usage(
    id UUID NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    model VARCHAR(256) NOT NULL,
    input_tokens INT NOT NULL,
    output_tokens INT NOT NULL,
    total_tokens INT NOT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Doc Store
############################################################
*/

CREATE TABLE folder(
    id BIGSERIAL,
    parent_id BIGINT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- index for maintaining unique folder names inside the same folder
CREATE UNIQUE INDEX idx_unique_folder_title_all
ON folder (customer_id, COALESCE(parent_id, -1), title);

CREATE TABLE document(
    id BIGSERIAL,
    parent_id BIGINT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    filename VARCHAR(1024) NOT NULL, -- human name of the document
    type VARCHAR(256) NOT NULL, -- txt, md, html, xlsx, etc.
    size_bytes BIGINT NOT NULL, -- size of the document in terms of bytes
    sha_256 CHAR(64) NOT NULL, -- a fingerprint of the document's contents
    validated BOOLEAN NOT NULL DEFAULT false, -- whether the object exists in datastore

    PRIMARY KEY (id),
    CONSTRAINT idx_unique_sha UNIQUE (customer_id, sha_256), -- no files can have same content anywhere

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- index for maintaining unique filenames inside the same folder
CREATE UNIQUE INDEX idx_unique_document_title_all
ON document (customer_id, COALESCE(parent_id, -1), filename);

-- vector objects that make up a document
CREATE TABLE document_vector(
    id BIGSERIAL,
    document_id BIGINT NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    vector_store_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    index INT NOT NULL, -- documents are chunked, so a large document will have multiple vector objects

    PRIMARY KEY (id),
    CONSTRAINT fk_vector_store FOREIGN KEY (vector_store_id, customer_id) REFERENCES vector_store(id, customer_id) ON DELETE CASCADE,
    CONSTRAINT fk_document_id_vector_store_id UNIQUE (document_id, vector_store_id, customer_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- asset information
CREATE TABLE asset_catalog(
    id BIGSERIAL,
    customer_id BIGINT,
    datastore TEXT NOT NULL, -- s3, etc.
    datastore_key UUID NOT NULL, -- remote id of the object. stored in assets/${uuid}
    filetype TEXT NOT NULL, -- what type of asset
    size_bytes BIGINT NOT NULL, -- how large the file is
    sha_256 CHAR(64) NOT NULL, -- fingerprint of the data

    PRIMARY KEY (id),
    CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Website
############################################################
*/

-- represents the website object for a user
CREATE TABLE website(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    protocol TEXT NOT NULL DEFAULT 'https',
    domain TEXT NOT NULL,
    blacklist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are disallowed
    whitelist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are allowed

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_website UNIQUE (customer_id, domain), -- websites are only allowed once

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- a page sourced from the sitemap of the website defined by the user
CREATE TABLE website_page(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    website_id BIGINT NOT NULL REFERENCES website(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    sha_256 CHAR(64) NOT NULL,
    is_valid BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_website_page UNIQUE (customer_id, website_id, url), -- pages are only allowed once

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- vectors associated with a website page's contents
CREATE TABLE website_page_vector(
    id BIGSERIAL,
    website_page_id BIGINT NOT NULL REFERENCES website_page(id) ON DELETE CASCADE,
    vector_store_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    index INT NOT NULL, -- data is chunked, so an index is required to sort the data

    PRIMARY KEY (id),
    CONSTRAINT fk_vector_store FOREIGN KEY (vector_store_id, customer_id) REFERENCES vector_store(id, customer_id) ON DELETE CASCADE,
    CONSTRAINT fk_website_page_id_vector_store_id UNIQUE (website_page_id, vector_store_id, customer_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Generated Content
############################################################
*/

-- for defining what types of content is supported
CREATE TABLE content_type(
    title VARCHAR(256) NOT NULL,
    parent TEXT NOT NULL DEFAULT 'Other', -- what list to sort this under

    PRIMARY KEY (title),
    CONSTRAINT unique_title_parent UNIQUE (title, parent),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- saved llm configurations for content creation
CREATE TABLE llm(
    id BIGSERIAL,
    customer_id BIGINT, -- when null the llm is a default for all customers
    model TEXT NOT NULL,
    temperature NUMERIC(1,2) NOT NULL,
    system_prompt TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- TODO : create default models for content

-- table to reference all generated content for the customer of all types
CREATE TABLE generation_library(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    content_type VARCHAR(256) NOT NULL,

    -- metadata
    draft BOOLEAN NOT NULL DEFAULT true,
    published BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_content_type FOREIGN KEY (content_type) REFERENCES content_type(title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Blog Post
############################################################
*/

-- root blog that a customer can create
-- customers can create mutliple blogs
-- contains generation configuration inside as well
CREATE TABLE blog(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,

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
    llm_content_generation_default_id BIGINT DEFAULT NULL,
    llm_vector_summarization_default_id BIGINT DEFAULT NULL,
    llm_website_summarization_default_id BIGINT DEFAULT NULL,
    llm_proof_reading_default_id BIGINT DEFAULT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    website_id BIGINT NOT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    folder_id BIGINT NOT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_idea_id BIGINT DEFAULT NULL,
    blog_category_id BIGINT DEFAULT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_id BIGINT NOT NULL,

    -- data
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    additional_instructions TEXT NOT NULL,
    asset_id BIGINT DEFAULT NULL, -- from asset_catalog
    metadata JSON NOT NULL DEFAULT '{}',

    -- models when null, uses customer/defined default
    content_generation_model_id BIGINT DEFAULT NULL,
    vector_summarization_model_id BIGINT DEFAULT NULL,
    website_summarization_model_id BIGINT DEFAULT NULL,
    proof_reading_model_id BIGINT DEFAULT NULL,

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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_id BIGINT NOT NULL,
    blog_post_section_id BIGINT NOT NULL,
    website_page_id BIGINT NOT NULL,
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
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
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
CREATE TABLE blog_section_content(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_id BIGINT NOT NULL,
    blog_post_section_id BIGINT NOT NULL,

    content TEXT NOT NULL, -- raw content that the user can edit / give feedback for
    feedback TEXT NOT NULL DEFAULT '', -- feedback is ALWAYS used after the content in the conversation
    index INT NOT NULL, -- index of the conversation

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_blog_section_content_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_section_content_blog_id FOREIGN KEY
    (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_section_content_blog_post_id FOREIGN KEY
    (blog_post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_section_content_blog_post_section_id FOREIGN KEY
    (blog_post_section_id) REFERENCES blog_post_section(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- tags on a blog post
-- multiple tags per blog post
CREATE TABLE blog_post_tag(
    id BIGSERIAL,
    customer_id BIGINT NOT NULL,
    blog_id BIGINT NOT NULL,
    blog_post_id BIGINT NOT NULL,

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

/*
############################################################
USERS
############################################################
*/

CREATE ROLE schema_spy LOGIN PASSWORD 'schema_spy';
GRANT CONNECT ON DATABASE aicontent TO schema_spy;
GRANT USAGE ON SCHEMA public TO schema_spy;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO schema_spy;

/*
############################################################
CUSTOM FUNCTION AND TRIGGERS
############################################################
*/

--
-- deleting old vector records when the joining table gets deleted
CREATE OR REPLACE FUNCTION delete_vector_if_unreferenced()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if there are no more references in document_vector
    IF (SELECT COUNT(*) FROM document_vector WHERE vector_store_id = OLD.vector_store_id) = 0 THEN
        -- Check if there are no more references in website_page_vector
        IF (SELECT COUNT(*) FROM website_page_vector WHERE vector_store_id = OLD.vector_store_id) = 0 THEN
            -- Delete from vector_store if there are no references
            DELETE FROM vector_store WHERE id = OLD.vector_store_id;
        END IF;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_delete_document_vector
AFTER DELETE ON document_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();

CREATE TRIGGER trg_after_delete_website_page_vector
AFTER DELETE ON website_page_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();

--
-- set llm default field to false when another record is set to be a default for the customer
CREATE OR REPLACE FUNCTION set_is_default_false()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the new or updated row is marked as default
    IF NEW.is_default THEN
        -- Update other rows
        UPDATE llm
        SET is_default = false
        WHERE customer_id = NEW.customer_id AND id != NEW.id AND is_default = true;
    END IF;
    -- Proceed with the insert or update
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_is_default_before_insert_or_update
BEFORE INSERT OR UPDATE ON llm
FOR EACH ROW EXECUTE FUNCTION set_is_default_false();
