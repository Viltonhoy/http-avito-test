create type operation_type as enum('deposit', 'withdrawal', 'transfer');

CREATE TABLE posting(
	id SERIAL PRIMARY KEY,
	account_id bigint NOT NULL,
	cb_journal operation_type NOT NULL,
	accounting_period text NOT NULL,
	amount bigint NOT NULL,
	date date NOT NULL,
	addressee bigint,
	description text 
);

create materialized view account_balances(
user_id, balance	
) as select
	account_id,
	sum(amount) 
from posting 
group by account_id
with no data;

create materialized view history_table(
account_id, cb_journal, amount, date, addressee, description
) as select
	account_id, cb_journal, amount, date, addressee, description 
from posting 
with no data;