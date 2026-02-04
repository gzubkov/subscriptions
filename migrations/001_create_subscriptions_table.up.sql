CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    start_date DATE NOT NULL,
    end_date DATE
);

CREATE INDEX idx_subscriptions_user_service_dates
    ON subscriptions (user_id, service_name, start_date, end_date);
CREATE UNIQUE INDEX idx_subscriptions_user_service_unique
    ON subscriptions (user_id, service_name);