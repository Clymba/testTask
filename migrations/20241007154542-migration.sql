-- +migrate Up

CREATE TABLE songs (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       group_name VARCHAR(255),
                       text TEXT,
                       genre VARCHAR(50),
                       date_added TIMESTAMP, -- добавляет поле для даты релиза
                       link VARCHAR(255)
)

-- +migrate Down

DROP TABLE IF EXISTS songs;

