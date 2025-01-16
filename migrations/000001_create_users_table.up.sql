CREATE TABLE IF NOT EXISTS users (
    id                UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name              VARCHAR(255) NOT NULL,
    username          VARCHAR(255) NOT NULL UNIQUE,
    email             VARCHAR(255) NOT NULL,
    password          VARCHAR(255) NOT NULL,
    gender            VARCHAR(10),
    dob               DATE,
    time_registration TIMESTAMP DEFAULT now(),
    avatar            TEXT UNIQUE
);

CREATE INDEX IF NOT EXISTS users_username_index ON users (username);