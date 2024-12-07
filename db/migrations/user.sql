-- 1. Создание пользователя
CREATE USER user_for_patefon WITH PASSWORD 'patefon';

-- 2. Предоставление прав
-- Замена "your_database" на имя вашей базы данных
GRANT CONNECT ON DATABASE patefon TO user_for_patefon; -- Доступ на подключение к бд.
GRANT USAGE ON SCHEMA public TO user_for_patefon;           -- Доступ к схеме (по умолчанию public)

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO user_for_patefon --Доступ на чтение и запись всех таблиц
