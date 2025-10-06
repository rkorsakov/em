CREATE TABLE IF NOT EXISTS subscriptions
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price        INTEGER      NOT NULL CHECK (price >= 0),
    user_id      VARCHAR(36)         NOT NULL,
    start_date   VARCHAR(7)   NOT NULL,
    end_date     VARCHAR(7),
    created_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions (service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_dates ON subscriptions (start_date, end_date);