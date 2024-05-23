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