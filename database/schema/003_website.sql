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
    path TEXT NOT NULL DEFAULT '',
    blacklist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are disallowed
    whitelist TEXT[] NOT NULL DEFAULT '{}', -- regex patterns that are allowed

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_website UNIQUE (customer_id, domain, path), -- websites are only allowed once

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
    metadata JSONB DEFAULT '{}',
    summary TEXT NOT NULL DEFAULT '',
    summary_sha_256 CHAR(64) NOT NULL DEFAULT '', -- fingerprint at the time the summary was taken

    vector_sha_256 CHAR(64) NOT NULL DEFAULT '', -- fingerprint when last vectorized

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
    metadata JSONB DEFAULT '{}',

    PRIMARY KEY (id),
    CONSTRAINT fk_vector_store FOREIGN KEY (vector_store_id, customer_id) REFERENCES vector_store(id, customer_id) ON DELETE CASCADE,
    CONSTRAINT fk_website_page_id_vector_store_id UNIQUE (website_page_id, vector_store_id, customer_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);