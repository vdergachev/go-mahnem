CREATE TABLE IF NOT EXISTS user_profile (
    user_profile_id SERIAL PRIMARY KEY,
    user_login VARCHAR (50) UNIQUE NOT NULL,
    user_name VARCHAR (50) NOT NULL,
    user_location_id INT REFERENCES user_location (user_location_id),
    user_language_id INT REFERENCES user_language (user_language_id),
    motto TEXT
);

-- add GRANT ...
