CREATE TABLE users
(
    id      BIGINT PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE transactions
(
    id          VARCHAR(100) PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    amount      BIGINT      NOT NULL,
    state       VARCHAR(10) NOT NULL,
    source_type VARCHAR(20) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);

INSERT INTO users (id, balance)
VALUES (1, 0),
       (2, 0),
       (3, 0)
ON CONFLICT (id) DO NOTHING;
