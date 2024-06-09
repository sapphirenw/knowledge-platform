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

-- list of available models and metadata about them
CREATE TABLE available_model(
    id VARCHAR(256) NOT NULL,
    provider TEXT NOT NULL, -- gemini, openai, etc.
    display_name TEXT NOT NULL,
    description TEXT NOT NULL,
    input_token_limit INT NOT NULL,
    output_token_limit INT NOT NULL,

    currency TEXT NOT NULL DEFAULT 'USD',
    input_cost_per_million_tokens NUMERIC(4,2) NOT NULL,
    output_cost_per_million_tokens NUMERIC(4,2) NOT NULL,

    depreciated_warning BOOLEAN NOT NULL DEFAULT false,
    is_depreciated BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),

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
    model VARCHAR(256) NOT NULL REFERENCES available_model(id) ON DELETE CASCADE, -- the model referenced
    temperature DOUBLE PRECISION NOT NULL,
    instructions TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    public BOOLEAN NOT NULL DEFAULT true, -- there will be internal models in some cases

    PRIMARY KEY (id),
    CONSTRAINT cnst_unqiue_llm_title UNIQUE
    (customer_id, title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_llm_title ON llm(title);

-- ties together conversation messages, can be used to seed an llm
CREATE TABLE conversation(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    title TEXT NOT NULL,
    conversation_type TEXT NOT NULL,
    system_message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',

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

    -- function call information
    tool_use_id TEXT NOT NULL DEFAULT '',
    tool_name TEXT NOT NULL DEFAULT '',
    tool_arguments JSONB NULL DEFAULT '{}',

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