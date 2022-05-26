create type operation_type as enum('deposit', 'withdrawal', 'transfer');
GO

CREATE TABLE posting(
	id SERIAL PRIMARY KEY,
	account_id bigint NOT NULL,
	cb_journal operation_type NOT NULL,
	accounting_period text NOT NULL,
	amount bigint NOT NULL,
	date date NOT NULL,
	addressee bigint 
);
GO

create materialized view account_balances(
user_id, balance	
) as select
	account_id,
	sum(amount) 
from posting 
group by account_id
with no data;
GO

create materialized view history_table(
account_id, cb_journal, amount, date, addressee
) as select
	account_id, cb_journal, amount, date, addressee 
from posting 
with no data;
Go