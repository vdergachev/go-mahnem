BEGIN;

CREATE TABLE IF NOT EXISTS user_location (
    user_location_id SERIAL PRIMARY KEY,
    country VARCHAR (50) NOT NULL,
    city VARCHAR (50) NOT NULL
);

GRANT SELECT, INSERT, UPDATE, DELETE ON public.user_location TO backend;

COMMIT;