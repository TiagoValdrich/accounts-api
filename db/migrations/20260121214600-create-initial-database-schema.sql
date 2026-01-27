
-- +migrate Up
CREATE TYPE transaction_operation_type AS ENUM (
    'normal_purchase',
    'installment_purchase',
    'withdrawal',
    'credit_voucher'
);

CREATE TABLE customer (
    id UUID PRIMARY KEY,
    document VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_document_unique UNIQUE (document)
);

CREATE TABLE customer_account (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_account_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customer(id)
);

CREATE INDEX idx_customer_account_customer_id ON customer_account(customer_id);

CREATE TABLE balance (
    id UUID PRIMARY KEY,
    customer_account_id UUID NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT balance_customer_account_id_fk FOREIGN KEY (customer_account_id) REFERENCES customer_account(id)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    customer_account_id UUID NOT NULL,
    operation_type transaction_operation_type NOT NULL,
    amount BIGINT NOT NULL,
    idempotency_key VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT transactions_customer_account_id_fk FOREIGN KEY (customer_account_id) REFERENCES customer_account(id),
    CONSTRAINT transactions_idempotency_key_unique UNIQUE (idempotency_key)
);

-- +migrate Down
DROP TABLE transactions;
DROP TABLE balance;
DROP TABLE customer_account;
DROP TABLE customer;
DROP TYPE transaction_operation_type;