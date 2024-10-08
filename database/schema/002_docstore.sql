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
    datastore_type TEXT NOT NULL DEFAULT 's3', -- s3, etc.
    datastore_id TEXT NOT NULL, -- id of the document in the remote datastore

    summary TEXT NOT NULL DEFAULT '',
    summary_sha_256 CHAR(64) NOT NULL DEFAULT '', -- fingerprint at the time the summary was taken

    vector_sha_256 CHAR(64) NOT NULL DEFAULT '', -- fingerprint when last vectorized

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
    metadata JSONB DEFAULT '{}',

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