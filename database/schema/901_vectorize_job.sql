-- job queue for vectorize jobs
CREATE TYPE vectorize_job_status AS ENUM ('waiting', 'in-progress', 'complete', 'error', 'unknown', 'rejected');
CREATE TABLE vectorize_job(
    id uuid NOT NULL DEFAULT uuid7(),
    customer_id uuid NOT NULL REFERENCES customer(id) ON DELETE CASCADE,

    documents BOOLEAN NOT NULL DEFAULT true,
    websites BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- a vectorize job item that stores the state of the job
CREATE TABLE vectorize_job_item(
    id uuid NOT NULL DEFAULT uuid7(),
    job_id uuid NOT NULL REFERENCES vectorize_job(id) ON DELETE CASCADE,
    
    status vectorize_job_status NOT NULL DEFAULT 'waiting',
    message TEXT NOT NULL DEFAULT 'Waiting ...',
    error TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);