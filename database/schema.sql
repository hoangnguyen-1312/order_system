CREATE TABLE user_account (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hashing VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    notes TEXT,
    seller_id SERIAL,
    price INTEGER,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES user_account (id) ON DELETE CASCADE
); 

CREATE TABLE purchase (
    id SERIAL PRIMARY KEY,
    buyer_id SERIAL,
    item_id SERIAL,
    price INTEGER,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (buyer_id) REFERENCES user_account (id),
    FOREIGN KEY (item_id) REFERENCES item (id)
);