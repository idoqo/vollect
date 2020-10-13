BEGIN;
CREATE TABLE vollect_tasks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    payload BYTEA NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

COMMIT;

