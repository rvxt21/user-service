\c users
CREATE EXTENSION citext;
CREATE DOMAIN domain_email AS citext
CHECK(
   VALUE ~ '^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$'
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY UNIQUE,
    name VARCHAR(100),
    lastname VARCHAR(1000),
    role VARCHAR(20) NOT NULL,
    email domain_email UNIQUE NOT NULL,
    password VARCHAR(200) NOT NULL,
    login VARCHAR(20) UNIQUE,
    status VARCHAR(10) NOT NULL DEFAULT 'active'
);

INSERT INTO users(email, login, password) VALUES ('anastasiia@gmail.com', 'anare', '123456Ql') ON CONFLICT (email) DO NOTHING;