-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    balance BIGINT NOT NULL
);

-- Insert initial data
INSERT INTO accounts (id, balance) VALUES
    (1, 1000),
    (2, 2000)
ON CONFLICT (id) DO NOTHING; 