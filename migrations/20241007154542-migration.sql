-- +migrate Up

CREATE TABLE songs (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- ID песни, генерируется автоматически
                       group_name VARCHAR(255),                        -- Название группы
                       text TEXT,                                      -- Текст песни
                       genre VARCHAR(50),                              -- Жанр песни
                       date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата добавления (по умолчанию текущее время)
                       link VARCHAR(255)                               -- Ссылка на песню
);

-- +migrate Down

DROP TABLE IF EXISTS songs;