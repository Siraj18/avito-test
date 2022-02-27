package postgresdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/siraj18/avito-test/internal/models"
)

var ErrorUserNotFound = fmt.Errorf("user not found")
var ErrorNotEnoughMoney = fmt.Errorf("not enough money")

const (
	operationAddMoney      = "adding money"
	operationWithdrawMoney = "withdrawal of money"
	operationTransferMoney = "transfer money"
)

func (rep *SqlRepository) createUserBalance(uid string, tx *sql.Tx) error {
	_, err := tx.Exec(addUserSql, uid)

	return err
}

func (rep *SqlRepository) GetBalance(uid string) (*models.User, error) {
	var user models.User

	if err := rep.db.Get(&user, getUserBalanceSql, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error when get user balance: %w", err)
	}

	return &user, nil

}

func (rep *SqlRepository) ChangeBalance(uid string, money float64) (*models.User, error) {
	tx, err := rep.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user models.User

	if err := tx.QueryRow(getUserSql, uid).Scan(&user.Id); err != nil {
		if err == sql.ErrNoRows {
			err = rep.createUserBalance(uid, tx)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}
	row := tx.QueryRow(updateUserBalanceSql, uid, money)
	if err = row.Scan(&user.Id, &user.Balance); err != nil {
		//TODO проверка ошибок
		if strings.Contains(err.Error(), "users_balance_check") {
			return nil, ErrorNotEnoughMoney
		}
		return nil, err
	}

	if money > 0 {
		if err = rep.addTransaction(&uid, nil, operationAddMoney, money, tx); err != nil {
			return nil, err
		}
	} else {
		if err = rep.addTransaction(nil, &uid, operationWithdrawMoney, money, tx); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (rep *SqlRepository) TransferBalance(fromUid string, toUid string, money float64) error {
	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(updateUserBalanceSql, toUid, money)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotFound
		}

		return err
	}

	_, err = tx.Exec(updateUserBalanceSql, fromUid, -money)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotFound
		} else if strings.Contains(err.Error(), "users_balance_check") {
			return ErrorNotEnoughMoney
		}

		return err
	}

	if err = rep.addTransaction(&toUid, &fromUid, operationTransferMoney, money, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
