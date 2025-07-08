-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный ID платежа
    user_id Bigserial NOT NULL,                         -- Пользователь
    role TEXT NOT NULL,                            -- Роль: student или parent
    invoice_id TEXT UNIQUE NOT NULL,               -- ID от AirbaPay
    amount NUMERIC NOT NULL,                       -- Сумма
    quantity INTEGER NOT NULL,                     -- Кол-во товаров
    status TEXT NOT NULL,                          -- Статус (например, "success", "failed")
    created_at TIMESTAMP DEFAULT NOW()             -- Дата создания
);
