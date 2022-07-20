create type operation_type as enum('deposit', 'withdrawal', 'transfer');

CREATE TABLE posting(
	id SERIAL PRIMARY KEY,
	account_id bigint NOT NULL,
	cb_journal operation_type NOT NULL,
	accounting_period date NOT NULL,
	amount bigint NOT NULL,
	date timestamp with time zone NOT NULL,
	addressee bigint,
	description text 
);

CREATE TABLE account_balances(
	account_id bigint NOT NULL unique,
	balances bigint 
);