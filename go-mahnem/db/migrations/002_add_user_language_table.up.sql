CREATE TABLE IF NOT EXISTS user_language (
    user_language_id SERIAL PRIMARY KEY,
    language_name VARCHAR(50) UNIQUE NOT NULL
);

-- add GRANT ...
