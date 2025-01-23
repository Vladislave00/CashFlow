CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255)
);

CREATE TABLE accountings (
    accounting_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    owner_id INT REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE money_values (
	value_id SERIAL PRIMARY KEY,
	name VARCHAR(50),
	rub FLOAT
);

CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    accounting_id INT REFERENCES accountings(accounting_id) ON DELETE CASCADE,
    name VARCHAR(100),
    money_amount FLOAT,
    value_id INT REFERENCES money_values(value_id)
);

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,
    external_account_id INT REFERENCES accounts(account_id),
    money_amount FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_account_balances()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE accounts
    SET money_amount = money_amount - NEW.money_amount
    WHERE account_id = NEW.account_id;

    UPDATE accounts
    SET money_amount = money_amount + NEW.money_amount
    WHERE account_id = NEW.external_account_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_transaction_insert
AFTER INSERT ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_account_balances();

CREATE OR REPLACE FUNCTION create_external_account()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO accounts (accounting_id, name, money_amount, value_id)
    VALUES (NEW.accounting_id, 'Внешний счёт', 0, 1);  -- Устанавливаем начальную сумму 0 и значение NULL для value_id
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_accounting_insert
AFTER INSERT ON accountings
FOR EACH ROW
EXECUTE FUNCTION create_external_account();
