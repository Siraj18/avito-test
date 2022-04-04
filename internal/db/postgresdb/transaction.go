package postgresdb

import (
	"database/sql"
	"fmt"

	"github.com/siraj18/avito-test/internal/models"
)

var ErrorTransactionNotFound = fmt.Errorf("transaction not found")

func (rep *SqlRepository) addTransaction(toId, fromId *string, operation string, money float64, tx *sql.Tx) error {
	_, err := tx.Exec(addTransactionsSql, toId, fromId, money, operation)

	return err
}

func (rep *SqlRepository) GetTransaction(id string) (*models.Transaction, error) {
	var transaction models.Transaction

	err := rep.db.Get(&transaction, getTransactionSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorTransactionNotFound
		}

		return nil, err
	}

	return &transaction, nil
}

func (rep *SqlRepository) GetAllTransactions(id string, sortType string, limit int, page int) (*[]models.Transaction, error) {
	finalSql := ""

	switch sortType {
	case "date_asc":
		finalSql = fmt.Sprintf(getAllTransactionsSql, "ORDER BY created_at ASC")
	case "date_desc":
		finalSql = fmt.Sprintf(getAllTransactionsSql, "ORDER BY created_at DESC")
	case "money_asc":
		finalSql = fmt.Sprintf(getAllTransactionsSql, "ORDER BY money ASC")
	case "money_desc":
		finalSql = fmt.Sprintf(getAllTransactionsSql, "ORDER BY money DESC")
	default:
		finalSql = getAllTransactionsSql
	}

	transactions := []models.Transaction{}

	offset := (page - 1) * limit

	err := rep.db.Select(&transactions, finalSql, id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &transactions, nil

}
