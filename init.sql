-- Enable UUID extension (built-in for Postgres 13+)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user' NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 2. Refresh Tokens Table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 3. JWT Denylist Table
CREATE TABLE IF NOT EXISTS jwt_denylist (
    jti VARCHAR(255) PRIMARY KEY,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 4. Products Table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    icode INTEGER NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    batch_no INTEGER NOT NULL,
    mrp DECIMAL(10, 2) NOT NULL,
    barcode VARCHAR(100) NOT NULL,
    UNIQUE (icode, batch_no)
);

-- Migration (run on existing databases):
-- ALTER TABLE products DROP COLUMN IF EXISTS created_at;

-- Index on barcode to optimize product lookups
CREATE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);

-- 5. Inventory Logs Table
CREATE TABLE IF NOT EXISTS inventory_logs (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    updated BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Migration (run on existing databases):
-- ALTER TABLE inventory_logs ADD COLUMN IF NOT EXISTS updated BOOLEAN DEFAULT FALSE NOT NULL;
