-- Parent Tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) CHECK (role IN ('admin', 'customer', 'publisher')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE developers (
    id SERIAL PRIMARY KEY,
    developer_name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    genre_name VARCHAR(100) UNIQUE NOT NULL
);

-- User Profiles (1:1 with users)
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    customer_name VARCHAR(255) NOT NULL,
    current_balance NUMERIC(12, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE publishers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    publisher_name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    security_level VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Game Management
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    publisher_id INT REFERENCES publishers(id),
    developer_id INT REFERENCES developers(id),
    game_name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock_level INT DEFAULT 0,
    release_date DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Mapping and History
CREATE TABLE game_genres (
    game_id INT REFERENCES games(id) ON DELETE CASCADE,
    genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (game_id, genre_id)
);

CREATE TABLE game_quantity_history (
    id SERIAL PRIMARY KEY,
    game_id INT REFERENCES games(id),
    change_amount INT NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Financials and Library
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id),
    game_id INT REFERENCES games(id),
    order_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    qty INT NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ledger (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id),
    order_id INT REFERENCES orders(id),
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(10) CHECK (type IN ('credit', 'debit')),
    amount NUMERIC(12, 2) NOT NULL
);

CREATE TABLE customer_game_library (
    customer_id INT REFERENCES customers(id),
    game_id INT REFERENCES games(id),
    purchase_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (customer_id, game_id)
);