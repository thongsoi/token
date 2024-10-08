CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id INT REFERENCES users(id)
);
