BEGIN;

CREATE TABLE IF NOT EXISTS user_to_language (
    user_profile_id INT NOT NULL,
    user_language_id INT NOT NULL,
    FOREIGN KEY (user_profile_id) REFERENCES user_profile (user_profile_id),
    FOREIGN KEY (user_language_id) REFERENCES user_language (user_language_id)
);

CREATE UNIQUE INDEX ux_user_profile_to_user_language
    ON user_to_language (user_profile_id, user_language_id);
    
GRANT SELECT, INSERT, UPDATE, DELETE ON public.user_to_language TO backend;

COMMIT;

