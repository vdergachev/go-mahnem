BEGIN;

CREATE TABLE IF NOT EXISTS user_profile (
    user_profile_id SERIAL PRIMARY KEY,
    user_login VARCHAR (60) UNIQUE NOT NULL,
    user_name VARCHAR (60) NOT NULL,
    user_location_id INT NOT NULL,
    motto TEXT,
    instagram_login TEXT,
    created_date TIMESTAMP NOT NULL,
    last_updated_date TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (user_location_id) REFERENCES user_location (user_location_id) ON DELETE CASCADE
);

GRANT SELECT, INSERT, UPDATE, DELETE ON public.user_profile TO backend;

COMMIT;
