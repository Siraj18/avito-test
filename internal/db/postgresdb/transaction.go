package postgresdb

import "database/sql"

func (rep *SqlRepository) addTransaction(toId, fromId *string, operation string, money float64, tx *sql.Tx) error {
	_, err := tx.Exec(addTransactionsSql, toId, fromId, money, operation)

	return err
}
