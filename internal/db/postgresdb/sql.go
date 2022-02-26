package postgresdb

const initSchema = `
				DROP TABLE IF EXISTS transactions;
				DROP TABLE IF EXISTS users;
				
				CREATE TABLE IF NOT EXISTS users 
				(
					id  	UUID PRIMARY KEY,
					balance DECIMAL(10, 2) DEFAULT 0
				);
				CREATE TABLE IF NOT EXISTS transactions
				(
					id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					to_id      UUID REFERENCES users(id),
					from_id    UUID REFERENCES users(id),
					money      DECIMAL(10, 2) NOT NULL,
					operation  TEXT NOT NULL,
					created_at TIMESTAMP DEFAULT now()
				);
				CREATE INDEX ON transactions (to_id);
				CREATE INDEX ON transactions (from_id);
`

const addUserSql = `
				INSERT INTO users VALUES ($1, 0);
`

const updateUserBalanceSql = `
				UPDATE users SET balance=balance + $2
				WHERE id=$1;	
`

const getUserBalanceSql = `
				SELECT id, balance FROM users
				WHERE id=$1;
`

const getUserSql = `
				SELECT id FROM users
				WHERE id=$1;
`

const addTransactionsSql = `
				INSERT INTO transactions (to_id, from_id, money, operation)
				VALUES ($1, $2, $3, $4);
`
