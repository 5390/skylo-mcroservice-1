CREATE TABLE received_messages (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);