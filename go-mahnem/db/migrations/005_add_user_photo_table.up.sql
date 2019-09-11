BEGIN;

CREATE TABLE IF NOT EXISTS user_photo (
    user_photo_id SERIAL PRIMARY KEY,
    user_profile_id INT NOT NULL,
    url VARCHAR (200) NOT NULL,

    FOREIGN KEY (user_profile_id) REFERENCES user_profile (user_profile_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_user_photo
    ON user_photo (user_profile_id, url);

GRANT SELECT, INSERT, UPDATE, DELETE ON public.user_photo TO backend;

COMMIT;

