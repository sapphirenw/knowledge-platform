/*
############################################################
Beta Auth
############################################################
*/

-- Beta auth table to handle basic authentication for beta testers.
-- this does NOT maintain a relationship between customers and api keys,
-- and should be replaced semi-quickly
CREATE TABLE beta_api_key(
    id uuid NOT NULL DEFAULT uuid7(),
    name TEXT NOT NULL, -- the name associated with this api key
    expired BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);