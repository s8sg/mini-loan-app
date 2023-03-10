CREATE TABLE IF NOT EXISTS loans
(
    id          UUID PRIMARY KEY,
    customer_id VARCHAR NOT NULL,
    amount      NUMERIC NOT NULL,
    term        INT NOT NULL,
    status      VARCHAR NOT NULL,
    start_date  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customer_id_loans ON loans (customer_id);


CREATE TABLE IF NOT EXISTS repayments
(
    id          UUID PRIMARY KEY,
    num         INT NOT NULL,
    loan_id     UUID,
    amount      NUMERIC NOT NULL,
    status      VARCHAR NOT NULL,
    due_date    TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_customer_id_repayments ON loans (customer_id);

