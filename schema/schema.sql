CREATE TABLE customer(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    datastore VARCHAR(256) NOT NULL DEFAULT 's3', -- name of the datastore the user wants to store their documents

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE folder(
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE document(
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL REFERENCES folder(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    filename VARCHAR(1024) NOT NULL, -- human name of the document
    type VARCHAR(256) NOT NULL, -- txt, md, html, xlsx, etc.
    size_bytes BIGINT NOT NULL, -- size of the document in terms of bytes
    sha_256 CHAR(64) NOT NULL, -- acts as a fingerprint of the document to compare old vs new

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vector_store(
    id BIGSERIAL PRIMARY KEY,

    -- vector data
    raw TEXT NOT NULL,
    embeddings VECTOR(512) NOT NULL,

    -- metadata
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    document_id BIGINT NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    index INTEGER NOT NULL, -- documents are chunked, so the index of the chunk

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE token_usage(
    id UUID NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    model VARCHAR(256) NOT NULL,
    input_tokens INT NOT NULL,
    output_tokens INT NOT NULL,
    total_tokens INT NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);