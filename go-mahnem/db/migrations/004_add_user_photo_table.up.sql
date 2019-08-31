CREATE TABLE IF NOT EXISTS user_photo (
    user_profile_id INT NOT NULL,
    url VARCHAR (200) NOT NULL,
    FOREIGN KEY (user_profile_id) REFERENCES user_profile (user_profile_id)
);

-- add GRANT ...

