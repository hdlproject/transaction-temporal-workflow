CREATE TYPE transaction_status AS ENUM ('CREATED', 'PENDING', 'SUCCESS', 'FAILED');

CREATE TABLE IF NOT EXISTS product
(
    code  VARCHAR(50) PRIMARY KEY,
    name  VARCHAR(100) NOT NULL,
    price INTEGER      NOT NULL
);

CREATE TABLE IF NOT EXISTS "user"
(
    id      VARCHAR(50) PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    balance INTEGER      NOT NULL default 0,
    CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS user_balance_event
(
    id                 BIGSERIAL PRIMARY KEY,
    user_id            VARCHAR(50)        NOT NULL,
    balance            INTEGER            NOT NULL default 0,
    transaction_id     VARCHAR(100)       NOT NULL,
    transaction_status transaction_status NOT NULL,
    created_at         timestamptz        NOT NULL,
    is_published       BOOLEAN            NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user_balance_event_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE TABLE IF NOT EXISTS transaction
(
    id             BIGSERIAL PRIMARY KEY,
    transaction_id VARCHAR(100)       NOT NULL,
    status         transaction_status NOT NULL,
    amount         INTEGER            NOT NULL,
    product_code   VARCHAR(100)       NOT NULL,
    user_id        VARCHAR(100)       NOT NULL,
    created_at     timestamptz        NOT NULL,
    is_published   BOOLEAN            NOT NULL default FALSE,
    CONSTRAINT fk_transaction_product FOREIGN KEY (product_code) REFERENCES product (code),
    CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);
