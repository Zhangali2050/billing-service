-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,                    -- Уникальный числовой ID (1,2,3…)
    user_id BIGINT NOT NULL,                     -- ID пользователя
    role TEXT NOT NULL,                          -- Роль: student или parent
    invoice_id TEXT UNIQUE NOT NULL,             -- ID от AirbaPay
    amount NUMERIC NOT NULL,                     -- Сумма
    quantity INTEGER NOT NULL,                   -- Кол-во товаров
    status TEXT NOT NULL,                        -- Статус (например, "success", "failed")
    created_at TIMESTAMP DEFAULT NOW()           -- Дата создания
);
