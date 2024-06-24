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
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    raw TEXT NOT NULL, -- string utf-8 representation of the data 
    embeddings VECTOR(512) NOT NULL,
    content_type TEXT NOT NULL, -- the type of content. document, website_page, etc.
    object_id uuid NOT NULL, -- id that this object ties to
    object_parent_id uuid, -- id that the object's parent is tied to
    metadata JSONB DEFAULT '{}', -- includes AT LEAST two fields: object

    PRIMARY KEY (id, customer_id), -- customer_id needs to exist in the key for partitioning

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) PARTITION BY LIST(customer_id);
CREATE INDEX ON vector_store USING hnsw (embeddings vector_ip_ops);
CREATE INDEX idx_vector_store_object_id ON vector_store(object_id);
CREATE TABLE vector_store_default PARTITION OF vector_store DEFAULT; -- default

CREATE TYPE vectorize_status AS ENUM ('waiting', 'in-progress', 'complete', 'error', 'unknown', 'rejected');
CREATE TABLE vectorize_job(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id),
    status vectorize_status NOT NULL DEFAULT 'waiting',
    message TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);