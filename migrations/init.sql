-- Таблица ролей пользователей
CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('student', 'parent'))
);

-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student', 'parent')),
    amount INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS saved_cards (
    id SERIAL PRIMARY KEY,
    account_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- psql -U postgres -d dev_db -h localhost -c "\i migrations/init.sql"
