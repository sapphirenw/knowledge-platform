/*
############################################################
LinkedIn Post
############################################################
*/

-- a single linked in account for the customer
-- this can be a user's account or an organization's account
CREATE TABLE linkedin_account(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    -- keys
    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE linkedin_reference_website(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    linkedin_account_id uuid NOT NULL,
    website_id uuid NOT NULL,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_linkedin_reference_website_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_linkedin_reference_website_linkedin_account_id FOREIGN KEY
    (linkedin_account_id) REFERENCES linkedin_account(id) ON DELETE CASCADE,
    CONSTRAINT fk_linkedin_reference_website_website_id FOREIGN KEY
    (website_id) REFERENCES website(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE linkedin_reference_folder(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL,
    linkedin_account_id uuid NOT NULL,
    folder_id uuid NOT NULL,

    -- keys
    PRIMARY KEY (id),
    CONSTRAINT fk_linkedin_reference_folder_customer_id FOREIGN KEY
    (customer_id) REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT fk_linkedin_reference_folder_linkedin_account_id FOREIGN KEY
    (linkedin_account_id) REFERENCES linkedin_account(id) ON DELETE CASCADE,
    CONSTRAINT fk_linkedin_reference_folder_folder_id FOREIGN KEY
    (folder_id) REFERENCES folder(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);