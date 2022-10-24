create type operation_type as enum('deposit', 'withdrawal', 'transfer');

create type expenses_type as enum('reservation', 'unreservation');

CREATE TABLE posting(
	id BIGSERIAL PRIMARY KEY,
	account_id bigint NOT NULL,
	cb_journal operation_type NOT NULL,
	accounting_period date NOT NULL,
	amount bigint NOT NULL,
	date timestamp with time zone NOT NULL,
	addressee bigint,
	description text 
);

CREATE TABLE balances(
	balance bigint NOT NULL,
    account_id bigint unique,
    last_tx_id bigint NOT NULL
);

CREATE TABLE deferred_expenses(
	account_id bigint NOT NULL, 
	service_id bigint NOT NULL,
	order_id bigint NOT NULL,
	operation expenses_type NOT NULL, 
	price bigint NOT NULL,
	tx_id      bigint references posting (id),
	UNIQUE (operation, order_id)
);

CREATE TABLE consolidated_report(
	account_id bigint NOT NULL, 
	service_id bigint NOT NULL,
	order_id bigint NOT NULL unique,
	sum bigint NOT NULL,
	tx_id      bigint references posting (id)
);
