CREATE TABLE IF NOT EXISTS user
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    user_id     uuid,
    nickname    text,
    phone       text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);
