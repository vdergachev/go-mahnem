
CREATE TABLE IF NOT EXISTS user_location (
    user_location_id SERIAL PRIMARY KEY,
    country VARCHAR (50) NOT NULL,
    city VARCHAR (50) NOT NULL
);

-- add GRANT ...
