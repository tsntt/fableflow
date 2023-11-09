CREATE TABLE IF NOT EXISTS transfers (
    id uuid PRIMARY KEY DEFAULT  gen_random_uuid(),
    receiver uuid,
    sender uuid,
    amount NUMERIC NOT NULL,
    status VARCHAR(20) NOT NULL,
    scheduled timestamp DEFAULT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id uuid PRIMARY KEY DEFAULT  gen_random_uuid(),
    bank_id uuid NOT NULL,
    balance NUMERIC,
    created_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS banks (
    id uuid PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS tmp_banks (
    id uuid PRIMARY KEY DEFAULT  gen_random_uuid(),
    domain VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    hash  TEXT NOT NULL UNIQUE,
    created_at timestamp NOT NULL
);

ALTER TABLE transfers ADD FOREIGN KEY (receiver) REFERENCES accounts (id);
ALTER TABLE transfers ADD FOREIGN KEY (sender) REFERENCES accounts (id);
ALTER TABLE accounts ADD FOREIGN KEY (bank_id) REFERENCES banks (id);