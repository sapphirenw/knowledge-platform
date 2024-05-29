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
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- reference to the conversation that generated this idea if applicable
    conversation_id uuid NULL REFERENCES conversation(id) ON DELETE SET NULL,

    title TEXT NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,

    -- TODO -- add vectors for similarity search

    PRIMARY KEY (id),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);