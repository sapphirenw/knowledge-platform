CREATE EXTENSION IF NOT EXISTS vector;

/*
############################################################
CUSTOMER
############################################################
*/

CREATE TABLE customer(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    datastore VARCHAR(256) NOT NULL DEFAULT 's3', -- name of the datastore the user wants to store their documents

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- generic table to hold vectors for all sorts of data
CREATE TABLE vector_store(
    id BIGSERIAL PRIMARY KEY,
    raw TEXT NOT NULL, -- string utf-8 representation of the data 
    embeddings VECTOR(512) NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- track token usage for a customer across multiple different models
CREATE TABLE token_usage(
    id UUID NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    model VARCHAR(256) NOT NULL,
    input_tokens INT NOT NULL,
    output_tokens INT NOT NULL,
    total_tokens INT NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
DOC STORE
############################################################
*/

CREATE TABLE folder(
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- index for maintaining unique folder names inside the same folder
CREATE UNIQUE INDEX idx_unique_folder_title_all
ON folder (customer_id, COALESCE(parent_id, -1), title);

CREATE TABLE document(
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    filename VARCHAR(1024) NOT NULL, -- human name of the document
    type VARCHAR(256) NOT NULL, -- txt, md, html, xlsx, etc.
    size_bytes BIGINT NOT NULL, -- size of the document in terms of bytes
    sha_256 CHAR(64) NOT NULL, -- a fingerprint of the document's contents
    validated BOOLEAN NOT NULL DEFAULT false, -- whether the object exists in datastore

    CONSTRAINT idx_unique_sha UNIQUE (customer_id, sha_256), -- no files can have same content anywhere

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- index for maintaining unique filenames inside the same folder
CREATE UNIQUE INDEX idx_unique_document_title_all
ON document (customer_id, COALESCE(parent_id, -1), filename);

-- vector objects that make up a document
CREATE TABLE document_vector(
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    vector_store_id BIGINT NOT NULL REFERENCES vector_store(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    index INT NOT NULL, -- documents are chunked, so a large document will have multiple vector objects

    CONSTRAINT fk_document_id_vector_store_id UNIQUE (document_id, vector_store_id),

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Website Ingest
############################################################
*/

-- represents the website object for a user
CREATE TABLE website(
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    protocol TEXT NOT NULL DEFAULT 'https',
    domain TEXT NOT NULL,
    blacklist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are disallowed
    whitelist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are allowed

    CONSTRAINT cnst_unique_website UNIQUE (customer_id, domain), -- websites are only allowed once

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- a page sourced from the sitemap of the website defined by the user
CREATE TABLE website_page(
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    website_id BIGINT NOT NULL REFERENCES website(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    sha_256 CHAR(64) NOT NULL,
    is_valid BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT cnst_unique_website_page UNIQUE (customer_id, website_id, url), -- pages are only allowed once

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- vectors associated with a website page's contents
CREATE TABLE website_page_vector(
    id BIGSERIAL PRIMARY KEY,
    website_page_id BIGINT NOT NULL REFERENCES website_page(id) ON DELETE CASCADE,
    vector_store_id BIGINT NOT NULL REFERENCES vector_store(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    index INT NOT NULL, -- data is chunked, so an index is required to sort the data

    CONSTRAINT fk_website_page_id_vector_store_id UNIQUE (website_page_id, vector_store_id),

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);