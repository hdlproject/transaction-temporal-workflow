CREATE TABLE IF NOT EXISTS product
(
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS "user"
(
    id   VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TYPE transaction_status AS ENUM ('CREATED', 'PENDING', 'SUCCESS', 'FAILED');

CREATE TABLE IF NOT EXISTS transaction
(
    id             BIGSERIAL PRIMARY KEY,
    transaction_id VARCHAR(100)       NOT NULL,
    status         transaction_status NOT NULL,
    amount         INTEGER            NOT NULL,
    product_code   VARCHAR(100)       NOT NULL,
    user_id        VARCHAR(100)       NOT NULL,
    created_at     timestamptz        NOT NULL,
    CONSTRAINT fk_transaction_product FOREIGN KEY (product_code) REFERENCES product (code),
    CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);
