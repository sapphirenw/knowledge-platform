-- +goose Up
-- +goose StatementBegin

-- add column to separate whether a document is an asset and to vectorize or not
ALTER TABLE document ADD COLUMN is_asset BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE document ADD COLUMN vectorize BOOLEAN NOT NULL DEFAULT true;

--
-- Begin resume
--

CREATE TABLE resume(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    title TEXT NOT NULL,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- files associated with this resume
CREATE TABLE resume_document(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    document_id uuid NOT NULL REFERENCES document(id) ON DELETE CASCADE,
    is_resume BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_resume_file UNIQUE (resume_id, document_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- resume websites
CREATE TABLE resume_website(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    website_id uuid NOT NULL REFERENCES website(id) ON DELETE CASCADE,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_resume_website UNIQUE (resume_id, website_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- resume website pages. These are NOT linked to the website, just if individual pages
-- need to be attached they can
CREATE TABLE resume_website_page(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    website_page_id uuid NOT NULL REFERENCES website_page(id) ON DELETE CASCADE,

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_resume_website_page UNIQUE (resume_id, website_page_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resume_about(
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,

    name TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',

    -- links
    github TEXT NOT NULL DEFAULT '',
    linkedin TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (resume_id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resume_work_experience(
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    index INT NOT NULL DEFAULT 0,
    
    -- block fields
    company TEXT NOT NULL DEFAULT '',
    position TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN NOT NULL DEFAULT false,

    -- raw information is meant to store ALL information related to this, not formatted
    information TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (resume_id, index),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE resume_project(
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    index INT NOT NULL DEFAULT 0,

    -- block fields
    title TEXT NOT NULL DEFAULT '',
    subtitle TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    start_date TIMESTAMP,
    end_date TIMESTAMP,

    -- raw information is meant to store ALL information related to this, not formatted
    information TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (resume_id, index),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resume_education(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,

    -- block fields
    institution TEXT NOT NULL DEFAULT '',
    major TEXT NOT NULL DEFAULT '',
    level TEXT NOT NULL DEFAULT '',
    gpa NUMERIC(2,1) NOT NULL DEFAULT 4.0,
    location TEXT NOT NULL DEFAULT '',
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN NOT NULL DEFAULT false,

    -- raw information is meant to store ALL information related to this, not formatted
    information TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resume_skill(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,

    title TEXT NOT NULL DEFAULT '',
    items TEXT[] NOT NULL DEFAULT '{}',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--
-- Specific to each job posting
--

CREATE TYPE resume_application_status AS ENUM (
    'not-started',
    'in-progress',
    'generated',
    'applied',
    'heard-back',
    'interviewing',
    'job-offer',
    'accepted'
);
CREATE TABLE resume_application(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,

    title TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    company_site TEXT NOT NULL DEFAULT '',
    raw_text TEXT NOT NULL DEFAULT '',
    status resume_application_status NOT NULL DEFAULT 'not-started',

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- to track conversations that a user has with this job posting. This can be about
-- skill gaps and how to shrink them, information gaps when generating content tailored
-- for this purpose
CREATE TABLE resume_application_conversation(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_id uuid NOT NULL REFERENCES resume(id) ON DELETE CASCADE,
    resume_application_id uuid NOT NULL REFERENCES resume_application ON DELETE CASCADE,
    conversation_id uuid NOT NULL REFERENCES conversation(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resume_application_key_word(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_application_id uuid NOT NULL REFERENCES resume_application ON DELETE CASCADE,

    title TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (id),
    CONSTRAINT cnst_unique_resume_application_key_word UNIQUE (title),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- final generated resume
CREATE TABLE resume_application_resume(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_application_id uuid NOT NULL REFERENCES resume_application ON DELETE CASCADE,
    document_id uuid NOT NULL REFERENCES document(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- final generated cover letter
CREATE TABLE resume_application_cover_letter(
    id uuid NOT NULL DEFAULT uuid7(),
    resume_application_id uuid NOT NULL REFERENCES resume_application ON DELETE CASCADE,
    document_id uuid NOT NULL REFERENCES document(id) ON DELETE CASCADE,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resume_application_cover_letter;
DROP TABLE resume_application_resume;
DROP TABLE resume_application_key_word;
DROP TABLE resume_application_conversation;
DROP TABLE resume_application;
DROP TYPE resume_application_status;
DROP TABLE resume_skill;
DROP TABLE resume_education;
DROP TABLE resume_project;
DROP TABLE resume_work_experience;
DROP TABLE resume_about;
DROP TABLE resume_website_page;
DROP TABLE resume_website;
DROP TABLE resume_document;
DROP TABLE resume;

ALTER TABLE document DROP COLUMN vectorize;
ALTER TABLE document DROP COLUMN is_asset;
-- +goose StatementEnd
