CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              sender_user_id INT NOT NULL,
                              transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('deposit', 'withdraw', 'transfer')),
                              amount NUMERIC(20, 8) NOT NULL DEFAULT 0,
                              receiver_user_id INT NOT NULL,
                              created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);