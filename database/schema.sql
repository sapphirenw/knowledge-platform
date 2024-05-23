
/*
 * MIT License
 *
 * Copyright (c) 2023-2024 Fabio Lima
 * 
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 * 
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 * 
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */
 
/**
 * Returns a time-ordered UUID (UUIDv6).
 *
 * Referencies:
 * - https://github.com/uuid6/uuid6-ietf-draft
 * - https://github.com/ietf-wg-uuidrev/rfc4122bis
 *
 * MIT License.
 *
 * Tags: uuid guid uuid-generator guid-generator generator time order rfc4122 rfc-4122
 */
create or replace function uuid6() returns uuid as $$
declare
begin
	return uuid6(clock_timestamp());
end $$ language plpgsql;

create or replace function uuid6(p_timestamp timestamp with time zone) returns uuid as $$
declare

	v_time numeric := null;

	v_gregorian_t numeric := null;
	v_clock_sequence_and_node numeric := null;

	v_gregorian_t_hex_a varchar := null;
	v_gregorian_t_hex_b varchar := null;
	v_clock_sequence_and_node_hex varchar := null;

	v_output_bytes bytea := null;

	c_100ns_factor numeric := 10^7::numeric;
	
	c_epoch numeric := -12219292800::numeric; -- RFC-4122 epoch: '1582-10-15'
	c_version bit(64) := x'0000000000006000'; -- RFC-4122 version: b'0110...'
	c_variant bit(64) := x'8000000000000000'; -- RFC-4122 variant: b'10xx...'

begin

	v_time := extract(epoch from p_timestamp);

	v_gregorian_t := (v_time - c_epoch) * c_100ns_factor;
	v_clock_sequence_and_node := random()::numeric * 2^62::numeric;

	v_gregorian_t_hex_a := lpad(to_hex((div(v_gregorian_t, 2^12::numeric)::bigint)), 12, '0');
	v_gregorian_t_hex_b := lpad(to_hex((mod(v_gregorian_t, 2^12::numeric)::bigint::bit(64) | c_version)::bigint), 4, '0');
	v_clock_sequence_and_node_hex := lpad(to_hex((v_clock_sequence_and_node::bigint::bit(64) | c_variant)::bigint), 16, '0');

	v_output_bytes := decode(v_gregorian_t_hex_a || v_gregorian_t_hex_b  || v_clock_sequence_and_node_hex, 'hex');

	return encode(v_output_bytes, 'hex')::uuid;
	
end $$ language plpgsql;

-------------------------------------------------------------------
-- EXAMPLE:
-------------------------------------------------------------------
-- 
-- select uuid6() uuid, clock_timestamp()-statement_timestamp() time_taken;
--
-- |uuid                                  |time_taken        |
-- |--------------------------------------|------------------|
-- |1eeca632-cf2a-65e0-85f3-151064c2409d  |00:00:00.000108   |
-- 

-------------------------------------------------------------------
-- EXAMPLE: generate a list
-------------------------------------------------------------------
-- 
-- with x as (select clock_timestamp() as t from generate_series(1, 10))
-- select uuid6(x.t) uuid, x.t::text ts from x;
-- 
-- |uuid                                |ts                           |
-- |------------------------------------|-----------------------------|
-- |1eeca634-f783-63f0-9988-48906d79f782|2024-02-13 08:30:37.891480-03|
-- |1eeca634-f783-6c24-97af-605238f4c3d0|2024-02-13 08:30:37.891691-03|
-- |1eeca634-f783-6e7c-9c2e-624f24b87738|2024-02-13 08:30:37.891754-03|
-- |1eeca634-f784-6070-a67b-4fc6659143e7|2024-02-13 08:30:37.891800-03|
-- |1eeca634-f784-6200-befd-0e20be5b0087|2024-02-13 08:30:37.891842-03|
-- |1eeca634-f784-6390-8f79-d4dacec1c3e0|2024-02-13 08:30:37.891881-03|
-- |1eeca634-f784-6520-8ee7-96091b017d4c|2024-02-13 08:30:37.891920-03|
-- |1eeca634-f784-66b0-a63e-c285d8a63e21|2024-02-13 08:30:37.891958-03|
-- |1eeca634-f784-6840-8c00-38659c4bf807|2024-02-13 08:30:37.891997-03|
-- |1eeca634-f784-69d0-b775-4bbfd45eb99e|2024-02-13 08:30:37.892036-03|
-- 

-------------------------------------------------------------------
-- FOR TEST: the expected result is an empty result set
-------------------------------------------------------------------
-- 
-- with t as (select uuid6() as id from generate_series(1, 1000))
-- select * from t where (id is null or id::text !~ '^[a-f0-9]{8}-[a-f0-9]{4}-6[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$');
--

/*
 * MIT License
 *
 * Copyright (c) 2023-2024 Fabio Lima
 * 
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 * 
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 * 
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */
 
/**
 * Returns a time-ordered with Unix Epoch UUID (UUIDv7).
 * 
 * Referencies:
 * - https://github.com/uuid6/uuid6-ietf-draft
 * - https://github.com/ietf-wg-uuidrev/rfc4122bis
 *
 * MIT License.
 *
 * Tags: uuid guid uuid-generator guid-generator generator time order rfc4122 rfc-4122
 */
create or replace function uuid7() returns uuid as $$
declare
begin
	return uuid7(clock_timestamp());
end $$ language plpgsql;

create or replace function uuid7(p_timestamp timestamp with time zone) returns uuid as $$
declare

	v_time numeric := null;

	v_unix_t numeric := null;
	v_rand_a numeric := null;
	v_rand_b numeric := null;

	v_unix_t_hex varchar := null;
	v_rand_a_hex varchar := null;
	v_rand_b_hex varchar := null;

	v_output_bytes bytea := null;

	c_milli_factor numeric := 10^3::numeric;  -- 1000
	c_micro_factor numeric := 10^6::numeric;  -- 1000000
	c_scale_factor numeric := 4.096::numeric; -- 4.0 * (1024 / 1000)
	
	c_version bit(64) := x'0000000000007000'; -- RFC-4122 version: b'0111...'
	c_variant bit(64) := x'8000000000000000'; -- RFC-4122 variant: b'10xx...'

begin

	v_time := extract(epoch from p_timestamp);

	v_unix_t := trunc(v_time * c_milli_factor);
	v_rand_a := ((v_time * c_micro_factor) - (v_unix_t * c_milli_factor)) * c_scale_factor;
	v_rand_b := random()::numeric * 2^62::numeric;

	v_unix_t_hex := lpad(to_hex(v_unix_t::bigint), 12, '0');
	v_rand_a_hex := lpad(to_hex((v_rand_a::bigint::bit(64) | c_version)::bigint), 4, '0');
	v_rand_b_hex := lpad(to_hex((v_rand_b::bigint::bit(64) | c_variant)::bigint), 16, '0');

	v_output_bytes := decode(v_unix_t_hex || v_rand_a_hex || v_rand_b_hex, 'hex');

	return encode(v_output_bytes, 'hex')::uuid;
	
end $$ language plpgsql;

-------------------------------------------------------------------
-- EXAMPLE:
-------------------------------------------------------------------
-- 
-- select uuid7() uuid, clock_timestamp()-statement_timestamp() time_taken;
--
-- |uuid                                  |time_taken        |
-- |--------------------------------------|------------------|
-- |018da240-e0db-72e1-86f5-345c2c240387  |00:00:00.000222   |
-- 

-------------------------------------------------------------------
-- EXAMPLE: generate a list
-------------------------------------------------------------------
-- 
-- with x as (select clock_timestamp() as t from generate_series(1, 1000))
-- select uuid7(x.t) uuid, x.t::text ts from x;
-- 
-- |uuid                                |ts                           |
-- |------------------------------------|-----------------------------|
-- |018da235-6271-70cd-a937-0bb7d22b801e|2024-02-13 08:23:44.113054-03|
-- |018da235-6271-7214-9188-1d3191883b5d|2024-02-13 08:23:44.113126-03|
-- |018da235-6271-723d-bebe-87f66085fad7|2024-02-13 08:23:44.113143-03|
-- |018da235-6271-728f-86ba-6e277d10c0a3|2024-02-13 08:23:44.113156-03|
-- |018da235-6271-72b8-9887-f31e4ca48020|2024-02-13 08:23:44.113168-03|
-- |018da235-6271-72e1-bbeb-8b686d0d4281|2024-02-13 08:23:44.113179-03|
-- |018da235-6271-730a-96a2-73275626f72a|2024-02-13 08:23:44.113190-03|
-- |018da235-6271-7333-8a5c-9d1ab89dc489|2024-02-13 08:23:44.113201-03|
-- |018da235-6271-735c-ba64-a42b55ad7d5c|2024-02-13 08:23:44.113212-03|
-- |018da235-6271-7385-a0fb-c65f5be24073|2024-02-13 08:23:44.113223-03|
--

-------------------------------------------------------------------
-- FOR TEST: the expected result is an empty result set
-------------------------------------------------------------------
-- 
-- with t as (select uuid7() as id from generate_series(1, 1000))
-- select * from t where (id is null or id::text !~ '^[a-f0-9]{8}-[a-f0-9]{4}-7[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$');
--


CREATE EXTENSION IF NOT EXISTS vector;

/*
############################################################
Base Schema
############################################################
*/

CREATE TABLE customer(
    id uuid NOT NULL DEFAULT uuid7(),
    name TEXT NOT NULL,
    datastore VARCHAR(256) NOT NULL DEFAULT 's3', -- name of the datastore the user wants to store their documents

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- generic table to hold vectors for all sorts of data
CREATE TABLE vector_store(
    id uuid NOT NULL DEFAULT uuid7(),
    raw TEXT NOT NULL, -- string utf-8 representation of the data 
    embeddings VECTOR(512) NOT NULL,
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    PRIMARY KEY (id, customer_id), -- customer_id needs to exist in the key for partitioning

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) PARTITION BY LIST(customer_id);
CREATE INDEX ON vector_store USING hnsw (embeddings vector_ip_ops);
CREATE TABLE vector_store_default PARTITION OF vector_store DEFAULT; -- default

/*
############################################################
Doc Store
############################################################
*/

CREATE TABLE folder(
    id uuid NOT NULL DEFAULT uuid7(),
    parent_id uuid NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_folder_title UNIQUE (customer_id, parent_id, title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE document(
    id uuid NOT NULL DEFAULT uuid7(),
    parent_id uuid NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    filename VARCHAR(1024) NOT NULL, -- human name of the document
    type VARCHAR(256) NOT NULL, -- txt, md, html, xlsx, etc.
    size_bytes BIGINT NOT NULL, -- size of the document in terms of bytes
    sha_256 CHAR(64) NOT NULL, -- a fingerprint of the document's contents
    validated BOOLEAN NOT NULL DEFAULT false, -- whether the object exists in datastore

    PRIMARY KEY (id),
    -- CONSTRAINT idx_unique_sha UNIQUE (customer_id, sha_256), -- no files can have same content anywhere
    CONSTRAINT cnst_unique_document_title UNIQUE (customer_id, parent_id, filename),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- vector objects that make up a document
CREATE TABLE document_vector(
    id uuid NOT NULL DEFAULT uuid7(),
    document_id uuid NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    vector_store_id uuid NOT NULL,
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    index INT NOT NULL, -- documents are chunked, so a large document will have multiple vector objects

    PRIMARY KEY (id),
    CONSTRAINT fk_vector_store FOREIGN KEY (vector_store_id, customer_id) REFERENCES vector_store(id, customer_id) ON DELETE CASCADE,
    CONSTRAINT fk_document_id_vector_store_id UNIQUE (document_id, vector_store_id, customer_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- asset information
CREATE TABLE asset_catalog(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid,
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
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
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
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    website_id uuid NOT NULL REFERENCES website(id) ON DELETE CASCADE,
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
    id uuid NOT NULL DEFAULT uuid7(),
    website_page_id uuid NOT NULL REFERENCES website_page(id) ON DELETE CASCADE,
    vector_store_id uuid NOT NULL,
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
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
-- the customer does not have access to this table
CREATE TABLE content_type(
    title VARCHAR(256) NOT NULL,
    parent TEXT NOT NULL DEFAULT 'Other', -- what list to sort this under

    PRIMARY KEY (title),
    CONSTRAINT cnst_unique_title_parent UNIQUE (title, parent),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- saved llm configurations for content creation
-- there are defaults for all customers to use, and customers can also save default configurations
CREATE TABLE llm(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid REFERENCES customer(id) ON DELETE CASCADE, -- when null the llm is a default for all customers

    title TEXT NOT NULL,
    color VARCHAR(7) DEFAULT NULL, -- #ffffff
    model TEXT NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    instructions TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unqiue_llm_title UNIQUE
    (customer_id, title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ties together conversation messages, can be used to seed an llm
CREATE TABLE conversation(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    title TEXT NOT NULL,
    conversation_type TEXT NOT NULL,
    metadata JSON DEFAULT '{}',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- a message in a conversation. stores a reference to the llm that generated it, in addition
-- to the static settings of the llm in case they change
CREATE TABLE conversation_message(
    id uuid NOT NULL DEFAULT uuid7(),
    conversation_id uuid NOT NULL REFERENCES conversation(id) ON DELETE CASCADE,
    llm_id uuid NULL REFERENCES llm(id) ON DELETE SET NULL,
    
    -- llm settings because these may change if the llm is updated
    model TEXT NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    instructions TEXT NOT NULL,

    -- conversation information
    role TEXT NOT NULL,
    message TEXT NOT NULL,
    index INT NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT cnst_conversation_message_unique UNIQUE
    (conversation_id, index),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- track token usage for a customer across multiple different models
CREATE TABLE token_usage(
    id UUID NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    -- optionally store a reference to the conversation that was coorelated with this token usage
    conversation_id uuid NULL REFERENCES conversation(id) ON DELETE SET NULL,

    model VARCHAR(256) NOT NULL,
    input_tokens INT NOT NULL,
    output_tokens INT NOT NULL,
    total_tokens INT NOT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
############################################################
Projects
############################################################
*/

-- Project for content generation. controls which documents are preferred default models
-- generation configs, etc.
-- Content is generated on a per-project basis.
-- Customers also have a root project which is automatically created
CREATE TABLE project(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    title TEXT NOT NULL,
    topic TEXT NOT NULL,
    idea_generation_model_id uuid DEFAULT NULL REFERENCES llm(id) ON DELETE SET NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);

-- folders that are tied to a project. When these are defined, content is preferentially pulled
-- from the documents owned by this folder
CREATE TABLE project_folder(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    folder_id uuid NOT NULL REFERENCES folder(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- websites tied to a project. Works the same as folders.
CREATE TABLE project_website(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    website_id uuid NOT NULL REFERENCES website(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- table that holds the base objects generated by each content type per project
-- acts as the reference to which project a post belongs to
CREATE TABLE project_library(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content_type VARCHAR(256) NOT NULL REFERENCES content_type(title),

    -- metadata
    draft BOOLEAN NOT NULL DEFAULT true,
    published BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ideas to use for content generation
CREATE TABLE project_idea(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    generation_batch_id uuid NOT NULL,

    -- reference to the conversation that generated this idea if applicable
    conversation_id uuid NULL REFERENCES conversation(id) ON DELETE SET NULL,

    title TEXT NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,

    -- TODO -- add vectors for similarity search

    PRIMARY KEY (id),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    project_library_id uuid NOT NULL REFERENCES project_library(id) ON DELETE CASCADE,

    additional_instructions TEXT NOT NULL DEFAULT '',
    
    title TEXT NOT NULL,
    asset_id uuid DEFAULT NULL REFERENCES asset_catalog(id) ON DELETE SET NULL,
    metadata JSON NOT NULL DEFAULT '{}',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE linkedin_post_content(
    id uuid NOT NULL DEFAULT uuid7(),
    linkedin_post_id uuid NOT NULL REFERENCES linkedin_post(id) ON DELETE CASCADE,

    content TEXT NOT NULL, -- raw content that the user can edit / give feedback for
    feedback TEXT NOT NULL DEFAULT '', -- feedback is ALWAYS used after the content in the conversation
    index INT NOT NULL, -- index of the conversation

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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
-- Set llm default field to false when another record is set to be a default for the customer
CREATE OR REPLACE FUNCTION set_is_default_false()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the new or updated row is marked as default
    IF NEW.is_default THEN
        -- Special handling for NULL customer_id (global default)
        IF NEW.customer_id IS NULL THEN
            -- Update other rows that are global defaults
            UPDATE llm
            SET is_default = false
            WHERE customer_id IS NULL AND id != NEW.id AND is_default = true;
        ELSE
            -- Update other rows for the same customer
            UPDATE llm
            SET is_default = false
            WHERE customer_id = NEW.customer_id AND id != NEW.id AND is_default = true;
        END IF;
    END IF;

    -- Proceed with the insert or update
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_is_default_before_insert_or_update
BEFORE INSERT OR UPDATE ON llm
FOR EACH ROW EXECUTE FUNCTION set_is_default_false();

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

/*
############################################################
USERS
############################################################
*/

CREATE ROLE schema_spy LOGIN PASSWORD 'schema_spy';
GRANT CONNECT ON DATABASE aicontent TO schema_spy;
GRANT USAGE ON SCHEMA public TO schema_spy;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO schema_spy;
